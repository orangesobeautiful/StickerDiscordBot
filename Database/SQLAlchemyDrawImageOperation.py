import sqlalchemy
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import func
from sqlalchemy.orm import sessionmaker
import math

_source_table_name = 'ImageSource'
_updated_folder_table_name = 'UpdatedFolder'
_warehouse_table_name = 'ImageWareHouse'

_Base = declarative_base()


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


class ImageSource(_Base):
    __tablename__ = _source_table_name

    folder_id = sqlalchemy.Column('folderid', sqlalchemy.String(512), primary_key=True, nullable=False)
    path = sqlalchemy.Column('path', sqlalchemy.Text, nullable=False)

    def __init__(self, folder_id, path):
        self.folder_id = folder_id
        self.path = path

    def __repr__(self):
        return "<ImageSource('{}', '{}')>".format(self.folder_id, self.path)

    def items(self):
        return [self.folder_id, self.path]


class ImageWareHouse(_Base):
    __tablename__ = _warehouse_table_name

    image_id = sqlalchemy.Column('imageid', sqlalchemy.String(512), primary_key=True, nullable=False)
    path = sqlalchemy.Column('path', sqlalchemy.Text, nullable=False)

    def __init__(self, image_id, path):
        self.image_id = image_id
        self.path = path

    def __repr__(self):
        return "<ImageWareHouse('{}', {})>".format(self.image_id, self.path)

    def items(self):
        return [self.image_id, self.path]


class UpdatedFolder(_Base):
    __tablename__ = _updated_folder_table_name

    folder_id = sqlalchemy.Column('folderid', sqlalchemy.String(512), primary_key=True, nullable=False)
    path = sqlalchemy.Column('path', sqlalchemy.Text, nullable=False)
    parent_folder_id = sqlalchemy.Column('parentfolderid', sqlalchemy.Text, nullable=False)

    def __init__(self, folder_id, parent_folder_id, path):
        self.folder_id = folder_id
        self.path = path
        self.parent_folder_id = parent_folder_id

    def __repr__(self):
        return "<UpdatedFolder('{}', '{}', {})>".format(self.folder_id, self.path, self.parent_folder_id)

    def items(self):
        return [self.folder_id, self.path, self.parent_folder_id]


class SQLAlchemyDrawImageOperation:
    _db_url = ''
    _session = None
    _engine = None

    def __init__(self, db_url: str):
        self._db_url = db_url
        self._init_db(db_url)

    def _init_db(self, db_url):
        print('DATABASE_URL=' + db_url)

        self._engine = create_engine(db_url, pool_recycle=2*3600, echo=False)
        self._session = sessionmaker(bind=self._engine)()

        self._create_tables()

    def _create_tables(self):
        # if table is not exist than create
        if not self._engine.dialect.has_table(self._engine, _source_table_name):
            ImageSource.metadata.create_all(self._engine)

        if not self._engine.dialect.has_table(self._engine, _warehouse_table_name):
            ImageWareHouse.metadata.create_all(self._engine)

        if not self._engine.dialect.has_table(self._engine, _updated_folder_table_name):
            UpdatedFolder.metadata.create_all(self._engine)

    # none test
    def add_image_source(self, folder_id, path):
        if folder_id and path:
            if not self.source_exist(folder_id):
                self._session.add(ImageSource(folder_id, path))
                self._session.commit()
            else:
                print('{0:s} 已經存在'.format(path))
                return False
            return True
        else:
            return False

    # none test
    def all_image_source(self):
        query_data = self._session.query(ImageSource.folder_id, ImageSource.path).order_by(ImageSource.path).all()

        return query_data

    # none test
    def source_exist(self, folder_id):
        query_data = self._session.query(ImageSource.path).filter(ImageSource.folder_id == folder_id).first()

        return query_data is not None

    # none test
    def delete_source(self, folder_id):
        if folder_id and self.source_exist(folder_id):
            fetch_num = self._session.query(ImageSource).filter(ImageSource.folder_id == folder_id).delete()
            self._session.commit()
            return fetch_num
        else:
            raise 'ID:{folder_id} is not exist'.format(folder_id)

    # none test
    def add_updated_folder(self, folder_id, path, parent_folder_id):
        if folder_id and path:
            if not self.source_exist(folder_id):
                self._session.add(UpdatedFolder(folder_id, path, parent_folder_id))
                self._session.commit()
            else:
                print('{0:s} 已經存在'.format(path))
                return False
            return True
        else:
            return False

    # none test
    def all_updated_folders(self, parent_folder=None):
        if parent_folder:
            query_data = self._session.query(UpdatedFolder.folder_id, UpdatedFolder.path, UpdatedFolder.parent_folder_id).\
                filter(UpdatedFolder.parent_folder_id == parent_folder).order_by(UpdatedFolder.folder_id).all()
        else:
            query_data = self._session.query(UpdatedFolder.folder_id, UpdatedFolder.path, UpdatedFolder.parent_folder_id).\
                order_by(UpdatedFolder.folder_id).all()
        return query_data

    # none test
    def delete_updated_folders(self):
        fetch_num = self._session.query(UpdatedFolder).delete()
        self._session.commit()

        return fetch_num

    # none test
    def add_images(self, image_id, path):
        if image_id and path:
            self._session.add(ImageWareHouse(image_id, path))
            self._session.commit()
            return True
        else:
            return False

    # none test
    def all_images(self):
        query_data = self._session.query(ImageWareHouse.image_id, ImageWareHouse.path).order_by(ImageWareHouse.path).all()

        return query_data

    # none test
    def delete_all_image(self):
        fetch_num = self._session.query(ImageWareHouse).delete()
        self._session.commit()
        return fetch_num

    # none test
    def get_rand_image(self):
        query_data = self._session.query(ImageWareHouse.image_id).order_by(func.random()).first()
        return query_data


if __name__ == '__main__':
    testDB = SQLAlchemyDrawImageOperation('mysql+pymysql://test:1234@localhost/our_bot')
    exit()



