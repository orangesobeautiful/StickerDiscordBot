import sqlalchemy
import os
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String, Boolean, DateTime
from sqlalchemy.sql import func
from sqlalchemy.orm import sessionmaker
import csv
import datetime

print(sqlalchemy.__version__)

db_url = os.environ['DATABASE_URL']
print(db_url)
engine = create_engine(db_url, echo=True)
Base = declarative_base()

exit()


class Sticker(Base):
    __tablename__ = 'sticker'

    id = Column(Integer, primary_key=True, nullable=False)
    name = Column('stickername', String, nullable=False)
    img_url = Column('imgurl', String, nullable=False)
    is_gif = Column('isgif', Boolean, nullable=False)
    latest_update_time = Column('latestupdatetime', DateTime(timezone=True), default=func.now(), nullable=False)

    def __init__(self, sticker_name, img_url, is_gif=False, latest_update_time=None):
        self.name = sticker_name
        self.img_url = img_url
        self.is_gif = is_gif
        self.latest_update_time = latest_update_time

    def __repr__(self):
        return "<Sticker('{}', '{}','{}', '{}', '{}')>".format(self.id, self.name, self.img_url, self.is_gif, self.latest_update_time)

    def items(self):
        return [self.id, self.name, self.img_url, self.is_gif, self.latest_update_time]


class BotInfo(Base):
    __tablename__ = 'botinfo'

    name = Column('name', String, primary_key=True, nullable=False)
    value = Column('value', String, nullable=False)

    def __init__(self, name, value):
        self.name = name
        self.value = value

    def __repr__(self):
        return "<BotInfo('{}', '{}')>".format(self.name, self.value)

    def items(self):
        return [self.name, self.value]


class ImageProxyGoogleDriver(Base):
    __tablename__ = 'imageproxygoogledriver'

    driver_id = Column('googledriverid', String, primary_key=True, nullable=False)
    local_file_name = Column('localfilename', String, nullable=False)
    latest_use_time = Column('latestusetime', String, nullable=False)

    def __init__(self, driver_id, file_name, latest_use_time):
        self.driver_id = driver_id
        self.local_file_name = file_name
        self.latest_use_time = latest_use_time

    def __repr__(self):
        return "<ImgaeProxyGoogleDriver('{}', '{}', '{}')>".format(self.name, self.value, self.latest_use_time)

    def items(self):
        return [self.driver_id, self.local_file_name, self.latest_use_time]


class ImageSource(Base):
    __tablename__ = 'imagesource'

    folder_id = Column('folderid', String, primary_key=True, nullable=False)
    path = Column('path', String, nullable=False)

    def __init__(self, folder_id, path):
        self.folder_id = folder_id
        self.path = path

    def __repr__(self):
        return "<ImageSource('{}', '{}')>".format(self.folder_id, self.path)

    def items(self):
        return [self.folder_id, self.path]


class ImageWareHouse(Base):
    __tablename__ = 'imagewarehouse'

    image_id = Column('imageid', String, primary_key=True, nullable=False)
    path = Column('path', String, nullable=False)

    def __init__(self, image_id, path):
        self.image_id = image_id
        self.path = path

    def __repr__(self):
        return "<ImageWareHouse('{}', {})>".format(self.image_id, self.path)

    def items(self):
        return [self.image_id, self.path]


class UpdatedFolder(Base):
    __tablename__ = 'updatefolder'

    folder_id = Column('folderid', String, primary_key=True, nullable=False)
    path = Column('path', String, nullable=False)
    parent_folder_id = Column('parentfolderid', String, nullable=False)

    def __init__(self, folder_id, parent_folder_id, path):
        self.folder_id = folder_id
        self.path = path
        self.parent_folder_id = parent_folder_id

    def __repr__(self):
        return "<UpdatedFolder('{}', '{}', {})>".format(self.folder_id, self.path, self.parent_folder_id)

    def items(self):
        return [self.folder_id, self.path, self.parent_folder_id]


Session = sessionmaker(bind=engine)
session = Session()

if not os.path.isdir('backup'):
    os.mkdir('backup')

backup_list = {
    'bot_info': BotInfo,
    'imageproxygooogledriver': ImageProxyGoogleDriver,
    'image_source': ImageSource,
    'image_warehouse': ImageWareHouse,
    'sticker': Sticker,
    'updated_folder': UpdatedFolder
}

for key in backup_list:
    with open('backup/' + key + '.csv', 'w', newline='', encoding='utf-8') as csvfile:
        writer = csv.writer(csvfile)

        for row in session.query(backup_list[key]).all():
            writer.writerow(row.items())