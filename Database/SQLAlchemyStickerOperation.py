import sqlalchemy
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import func
from sqlalchemy.orm import sessionmaker
from contextlib import contextmanager
import requests
from datetime import timezone
import math
import time
import os
import hashlib
import datetime
import csv
from CommonFunction import StickerCommon
import imghdr

project_dir = StickerCommon.project_dir
sticker_download_dir = StickerCommon.sticker_dir

_Base = declarative_base()

_sticker_table_name = 'sticker'
_bot_info_table_name = 'botinfo'


def imghdr_patch():

    def test_jpeg1(h, f):
        """JPEG data in JFIF format"""
        if b'JFIF' in h[:23]:
            return 'jpeg'

    JPEG_MARK = b'\xff\xd8\xff\xdb\x00C\x00\x08\x06\x06' \
                b'\x07\x06\x05\x08\x07\x07\x07\t\t\x08\n\x0c\x14\r\x0c\x0b\x0b\x0c\x19\x12\x13\x0f'

    def test_jpeg2(h, f):
        """JPEG with small header"""
        if len(h) >= 32 and 67 == h[5] and h[:32] == JPEG_MARK:
            return 'jpeg'

    def test_jpeg3(h, f):
        """JPEG data in JFIF or Exif format"""
        if h[6:10] in (b'JFIF', b'Exif') or h[:2] == b'\xff\xd8':
            return 'jpeg'

    imghdr.tests.append(test_jpeg1)
    imghdr.tests.append(test_jpeg2)
    imghdr.tests.append(test_jpeg3)


imghdr_patch()


def trans_url(img_url: str):
    if img_url[0: 25] == 'https://drive.google.com/':
        not_support_gd_url = False
        # google共用連結複製 or 直接連外網址
        if img_url[25: 25 + 5] == 'open?' or img_url[25: 25 + 3] == 'uc?':
            id_start = img_url.find('id=') + 3
            gd_file_id = img_url[id_start:]
        # google雲端上瀏覽
        elif img_url[25: 25 + 7] == 'file/d/':
            if img_url[-5:] == '/edit' or img_url[-5:] == '/view':
                gd_file_id = img_url[25 + 7:-5]
            elif '/' not in img_url[25 + 7:-5]:
                gd_file_id = img_url[25 + 7:]
            else:
                not_support_gd_url = True
        else:
            not_support_gd_url = True

        if not_support_gd_url:
            # await ctx.send('不支援的google driver網址格式')
            return False
        else:
            if not img_url == 'https://drive.google.com/uc?id=' + gd_file_id:
                img_url = 'https://drive.google.com/uc?id=' + gd_file_id

    return img_url


def download_image(url: str, filename: str):
    retry_times = 2
    success_download = False
    complete_filename = ''
    for _ in range(retry_times):
        try:
            response = requests.get(url, stream=True
                                    , headers={"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0"})
        except Exception:
            continue
        if response.ok:
            # detect image type
            r_iter = response.iter_content(4 * 1024)
            f_block = next(r_iter)
            img_type = imghdr.what(_, f_block)
            if img_type is None:
                success_download = False
                break

            complete_filename = filename + '.' + img_type

            # write image
            with open(os.path.join(sticker_download_dir, complete_filename), 'wb') as download_file:
                download_file.write(f_block)
                try:
                    for block in r_iter:
                        if not block:
                            break
                        download_file.write(block)
                except Exception:
                    continue

            success_download = True

        if success_download:
            break

        time.sleep(3)

    return success_download, complete_filename


class BotInfo(_Base):
    __tablename__ = _bot_info_table_name

    name = sqlalchemy.Column('name', sqlalchemy.String(512), primary_key=True, nullable=False)
    value = sqlalchemy.Column('value', sqlalchemy.Text, nullable=False)

    def __init__(self, name, value):
        self.name = name
        self.value = value

    def __repr__(self):
        return "<BotInfo('{}', '{}')>".format(self.name, self.value)

    def items(self):
        return [self.name, self.value]


class Sticker(_Base):
    __tablename__ = _sticker_table_name

    id = sqlalchemy.Column('id', sqlalchemy.Integer, primary_key=True, nullable=False, autoincrement=True)
    name = sqlalchemy.Column('stickername', sqlalchemy.Text, nullable=False)
    img_url = sqlalchemy.Column('imgurl', sqlalchemy.Text, nullable=False)
    local_save = sqlalchemy.Column('localsave', sqlalchemy.Text, nullable=False, default=func.now())
    is_gif = sqlalchemy.Column('isgif', sqlalchemy.Boolean, nullable=False)
    latest_update_time = sqlalchemy.Column('latestupdatetime', sqlalchemy.TIMESTAMP(timezone=True)
                                           , default=func.now(), nullable=False)

    def __init__(self, sticker_name, img_url, local_save='', is_gif=False, latest_update_time=None):
        self.name = sticker_name
        self.img_url = img_url
        self.local_save = local_save
        self.is_gif = is_gif
        self.latest_update_time = latest_update_time

    def __repr__(self):
        return "<Sticker('{}', '{}','{}', '{}', '{}')>".format(
            self.id, self.name, self.img_url, self.is_gif, self.latest_update_time)

    def items(self):
        return [self.id, self.name, self.img_url, self.is_gif, self.latest_update_time]


class SQLAlchemyStickerOperation:
    _db_url = ''
    _session = None
    _engine = None

    def __init__(self, db_url: str, save_image_local: bool, sticker_url: str= ""):
        self._db_url = db_url
        self.save_image_local = save_image_local
        self._init_db(db_url)
        self._sticker_url = sticker_url

    def _init_db(self, db_url):
        print('DATABASE_URL=' + db_url)

        self._engine = create_engine(db_url, pool_pre_ping=True, echo=False, pool_recycle=7200)
        self._session_maker = sessionmaker(bind=self._engine)

        self._create_tables()

    @contextmanager
    def _session_scope(self):
        """Provide a transactional scope around a series of operations."""
        session = self._session_maker()
        try:
            yield session
        except:
            session.rollback()
            raise
        finally:
            session.close()

    def _create_tables(self):
        insp = sqlalchemy.inspect(self._engine)
        # if table is not exist than create
        if not insp.has_table(_sticker_table_name, None):
            Sticker.metadata.create_all(self._engine)

        if not insp.has_table(_bot_info_table_name, None):
            BotInfo.metadata.create_all(self._engine)

    def get_bot_prefix(self):
        with self._session_scope() as session:
            query_data = session.query(BotInfo.value).filter(BotInfo.name == 'BotPrefix').first()
        if query_data is None:
            return None
        return query_data[0]

    def set_bot_prefix(self, prefix: str):
        with self._session_scope() as session:
            if self.get_bot_prefix() is None:
                session.add(BotInfo('BotPrefix', prefix))
            else:
                session.query(BotInfo).filter(BotInfo.name == 'BotPrefix').update({BotInfo.value: prefix})
            session.commit()

    def get_sticker_download_count(self):
        with self._session_scope() as session:
            query_data = session.query(BotInfo.value).filter(BotInfo.name == 'StickerDownloadCount').first()
            if query_data is None:
                session.add(BotInfo('StickerDownloadCount', 0))
                return 0
        return int(query_data[0])

    def set_sticker_download_count(self, count: int):
        with self._session_scope() as session:
            session.query(BotInfo).filter(BotInfo.name == 'StickerDownloadCount').update({BotInfo.value: str(count)})
            session.commit()

    def get_all_sn_list(self):
        with self._session_scope() as session:
            query_data = session.query(Sticker.name).distinct().all()
        res_list = list()
        for tup_ele in query_data:
            res_list.append(tup_ele[0])

        return res_list

    def get_sticker_random(self, sticker_name: str):
        with self._session_scope() as session:
            query_data = session.query(Sticker.img_url, Sticker.local_save, Sticker.is_gif).filter(
                Sticker.name == sticker_name).order_by(func.random()).first()

        return query_data

    def get_sticker_all(self, sticker_name: str):
        with self._session_scope() as session:
            query_data = session.query(Sticker.id, Sticker.img_url, Sticker.is_gif).filter(
                Sticker.name == sticker_name).order_by(Sticker.id).all()
        return query_data

    # 給網頁版查詢資料使用 直接回傳json
    def get_sticker_group_by_name(self, start: int, num: int = 10, sort_by: str = 'name'):
        if sort_by == 'name':
            sort_col = Sticker.name

        with self._session_scope() as session:
            sticker_name_list = session.query(Sticker.name).group_by(Sticker.name).order_by(sort_col).limit(num).offset(
                start).all()

        return_list = list()
        for sticker_name in sticker_name_list:
            sticker_prop_list = self.get_sticker_all(sticker_name[0])
            sticker_dict_list = list()
            for sticker_prop in sticker_prop_list:
                sid = sticker_prop[0]
                img_url = sticker_prop[1]
                is_gif = sticker_prop[2]
                sticker_dict_list.append({'id': sid, 'url': img_url, 'gif': is_gif})

            return_list.append([sticker_name[0], sticker_dict_list])

        # use mariadb >= 10.5 func
        # with self._session_scope() as session:
        #   query_data = session.query(Sticker.name, func.json_array_agg(func.json_object('id', Sticker.id,
        #                                                                                'url', Sticker.img_url,
        #                                                                                'gif', Sticker.is_gif))). \
        #    group_by(Sticker.name).order_by(sort_col).limit(num).offset(start).all()

        # use postgresql func
        # with self._session_scope() as session:
        #   query_data = session.query(Sticker.name, func.json_agg(func.json_build_object('id', Sticker.id,
        #                                                                                'url', Sticker.img_url,
        #                                                                                'gif', Sticker.is_gif))). \
        #    group_by(Sticker.name).order_by(sort_col).limit(num).offset(start).all()

        return return_list

    # 根據單頁顯示數量計算總頁數
    def max_page(self, num: int):
        if type(num) == str:
            num = int(num)
        return math.ceil(len(self.get_all_sn_list())/num)

    # 查詢同名貼圖是否存在
    def is_sticker_name_exist(self, sticker_name: str):
        with self._session_scope() as session:
            query_data = session.query(Sticker.id).filter(Sticker.name == sticker_name).first()
        return query_data is not None

    # 查詢同名同網址的貼圖是否存在
    def is_sticker_exist(self, sticker_name: str, img_url: str):
        with self._session_scope() as session:
            query_data = session.query(Sticker.id).filter(Sticker.name == sticker_name, Sticker.img_url == img_url).first()
        return query_data is not None

    # 判斷貼圖網址是否一樣(根據ID)
    def is_sticker_equal(self, sticker_id: str, img_url: str):
        with self._session_scope() as session:
            query_data = session.query(Sticker.name, Sticker.img_url).filter(Sticker.id == sticker_id).first()
        sticker_name = query_data[0]
        orgn_url = query_data[1]

        return orgn_url == img_url, sticker_name

    # 檢查本地圖片是否存在 找不到檔案則重新下載
    def check_local_save(self):
        with self._session_scope() as session:
            query_data = session.query(Sticker.id, Sticker.img_url, Sticker.local_save).all()
            download_failed_list = list()
            for sticker_data in query_data:
                s_id = sticker_data[0]
                img_url = sticker_data[1]
                local_save = sticker_data[2]

                if local_save == '':
                    s_count = int(self.get_sticker_download_count())
                    while True:
                        s_count += 1
                        hash_name = hashlib.md5(str(s_count).encode(encoding='utf-8')).hexdigest()
                        save_path = os.path.join(sticker_download_dir, hash_name)
                        if not os.path.isfile(save_path):
                            break
                    dl_success, complete_filename = download_image(img_url, hash_name)
                    if dl_success:
                        self.set_sticker_download_count(s_count)
                        local_save = complete_filename
                        session.query(Sticker).filter(Sticker.id == s_id).update({Sticker.local_save: local_save})
                    else:
                        download_failed_list.append((s_id, img_url))
                else:
                    save_path = os.path.join(sticker_download_dir, local_save)
                    if os.path.isfile(save_path):
                        # 本地圖片依然存在 檢查原網址是否可用 如果不可用則替換成本地圖片
                        if img_url != self.sticker_url + 'sticker-image/' + local_save:
                            if not download_image(img_url, save_path):
                                session.query(Sticker).filter(Sticker.id == s_id).update(
                                    {Sticker.img_url: self.sticker_url + 'sticker-image/' + local_save})
                    else:
                        if not download_image(img_url, save_path):
                            download_failed_list.append((s_id, img_url))

            session.commit()

        return download_failed_list

    def add_sticker(self, add_list: list):
        """
        err 1: not support url
        err 2: has equal sticker
        err 3: download failed
        """
        with self._session_scope() as session:
            no_add_list = list()
            for add_info in add_list:
                sticker_name: str = add_info['sn']
                img_url: str = add_info['url']
                is_gif = add_info['is_gif']
                if type(is_gif) == str:
                    if is_gif.lower() == 'false':
                        is_gif = False
                    elif is_gif.lower() == 'true':
                        is_gif = True

                if img_url[:4] == 'http':
                    af_url = trans_url(img_url)
                    if af_url:
                        if self.is_sticker_exist(sticker_name, img_url):
                            no_add_list.append({'sn': sticker_name, 'url': img_url, 'err': 2})
                        else:
                            local_save = ''
                            if self.save_image_local:
                                s_count = int(self.get_sticker_download_count())
                                while True:
                                    s_count += 1
                                    hash_name = hashlib.md5(str(s_count).encode(encoding='utf-8')).hexdigest()
                                    save_path = os.path.join(sticker_download_dir, hash_name)
                                    if not os.path.isfile(save_path):
                                        break
                                dl_success, complete_filename = download_image(img_url, hash_name)
                                if dl_success:
                                    self.set_sticker_download_count(s_count)
                                    local_save = complete_filename

                            session.add(Sticker(sticker_name, af_url, local_save, is_gif))
                    else:
                        no_add_list.append({'sn': sticker_name, 'url': img_url, 'err': 1})
                else:
                    no_add_list.append({'sn': sticker_name, 'url': img_url, 'err': 1})

            session.commit()
        return no_add_list

    def edit_sticker(self, edit_list: list):
        """
        err 1: not support url
        err 2: url no change
        err 3: has equal sticker
        err 4: error args
        """
        with self._session_scope() as session:
            no_change_list = list()
            for edit_info in edit_list:
                dy_update_dict = dict()
                sticker_id: str = str(edit_info['id'])
                img_url: str = ''
                if 'url' in edit_info:
                    img_url: str = edit_info['url']
                    dy_update_dict[Sticker.img_url] = img_url

                    local_save = ''
                    if self.save_image_local:
                        s_count = int(self.get_sticker_download_count())
                        while True:
                            s_count += 1
                            hash_name = hashlib.md5(str(s_count).encode(encoding='utf-8')).hexdigest()
                            save_path = os.path.join(sticker_download_dir, hash_name)
                            if not os.path.isfile(save_path):
                                break
                        dl_success, complete_filename = download_image(img_url, hash_name)
                        if dl_success:
                            self.set_sticker_download_count(s_count)
                            local_save = complete_filename
                        dy_update_dict[Sticker.local_save] = local_save

                if 'gif' in edit_info:
                    is_gif = edit_info['gif']
                    if type(is_gif) == str:
                        if is_gif.lower() == 'false':
                            is_gif = False
                        elif is_gif.lower() == 'true':
                            is_gif = True
                    dy_update_dict[Sticker.is_gif] = is_gif

                if len(dy_update_dict) == 0:
                    no_change_list.append({'id': sticker_id, 'img_url': img_url, 'err': 4})
                else:
                    if 'url' in edit_info:
                        af_url = trans_url(img_url)
                        if af_url:
                            equ_img, sticker_name = self.is_sticker_equal(sticker_id, img_url)
                            # print(equ_img, sticker_name)
                            if equ_img:
                                no_change_list.append({'id': sticker_id, 'img_url': img_url, 'err': 2})
                                continue
                            else:
                                if self.is_sticker_exist(sticker_name, img_url):
                                    no_change_list.append({'id': sticker_id, 'img_url': img_url, 'err': 3})
                                    continue
                        else:
                            no_change_list.append({'id': sticker_id, 'img_url': img_url, 'err': 1})
                            continue

                    session.query(Sticker).filter(Sticker.id == sticker_id).update(dy_update_dict)

            session.commit()
        return no_change_list

    def delete_sticker(self, id_list: list):
        with self._session_scope() as session:
            fetch_num = session.query(Sticker).filter(Sticker.id.in_(id_list)).delete(synchronize_session=False)
            session.commit()
        return fetch_num

    def delete_sticker_whole(self, sticker_name: str):
        with self._session_scope() as session:
            fetch_num = session.query(Sticker).filter(Sticker.name == sticker_name).delete()
            session.commit()
        return fetch_num


if __name__ == '__main__':
    op = SQLAlchemyStickerOperation('mysql+pymysql://test:1234@localhost/our_bot', True)
    op.check_local_save()

