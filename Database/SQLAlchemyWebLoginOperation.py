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
import random

_verification_code_length = 6

_Base = declarative_base()

_web_login_verification_table_name = 'webloginverification'
_web_user_info_table_name = 'webuserinfo'


class WebLoginVerification(_Base):
    __tablename__ = _web_login_verification_table_name

    code = sqlalchemy.Column('code', sqlalchemy.String(512), primary_key=True, nullable=False)
    expiration_time = sqlalchemy.Column('expirationtime', sqlalchemy.DateTime, nullable=False)
    user_id = sqlalchemy.Column('userid', sqlalchemy.Text, nullable=True)

    def __init__(self, code, expiration_time):
        self.code = code
        self.expiration_time = expiration_time

    def __repr__(self):
        return "<Web Login Verification('{}', '{}', '{}')>".format(self.code, self.expiration_time, self.user_id)


class WebUserInfo(_Base):
    __tablename__ = _web_user_info_table_name

    user_id = sqlalchemy.Column('userid', sqlalchemy.String(512), primary_key=True, nullable=False)
    name = sqlalchemy.Column('name', sqlalchemy.Text, nullable=False)
    avatar_url = sqlalchemy.Column('avatarurl', sqlalchemy.Text, nullable=False)

    def __init__(self, user_id, name, avatar_url):
        self.user_id = user_id
        self.name = name
        self.avatar_url = avatar_url

    def __repr__(self):
        return "<Web Login Verification('{}', '{}', '{}')>".format(self.user_id, self.name, self.avatar_url)


class SQLAlchemyWebLoginOperation:
    _db_url = ''
    _engine = None

    def __init__(self, db_url: str):
        self._db_url = db_url
        self._init_db(db_url)

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
        if not insp.has_table(_web_login_verification_table_name, None):
            WebLoginVerification.metadata.create_all(self._engine)
        if not insp.has_table(_web_user_info_table_name, None):
            WebUserInfo.metadata.create_all(self._engine)

    def is_code_exist(self, code: str):
        with self._session_scope() as session:
            query_data = session.query(WebLoginVerification.code, WebLoginVerification.expiration_time).filter(
                WebLoginVerification.code == code).first()

        if query_data is None:
            return False
        else:
            expiration_time = query_data[1]
            # 驗證時效過期
            if datetime.datetime.now() > expiration_time:
                self.delete_code(code)
                return False
            else:
                return True

    def generate_verification_code(self):
        code = random.randint(math.pow(10, _verification_code_length-1),
                              math.pow(10, _verification_code_length) - 1)
        while self.is_code_exist(code):
            code = random.randint(math.pow(10, _verification_code_length - 1),
                                  math.pow(10, _verification_code_length) - 1)
        expiration_time = datetime.datetime.now() + datetime.timedelta(minutes=5, seconds=5)
        with self._session_scope() as session:
            session.add(WebLoginVerification(str(code), expiration_time))
            session.commit()
        return code, expiration_time

    def user_id_sing_in(self, code: str, user_id: str):
        if self.is_code_exist(code):
            with self._session_scope() as session:
                tmp = session.query(WebLoginVerification).filter(
                    WebLoginVerification.code == code).update({WebLoginVerification.user_id: user_id},
                                                              synchronize_session=False)
                session.commit()
            return tmp
        else:
            return 0

    def check_verification_status(self, code: str):
        if self.is_code_exist(code):
            with self._session_scope() as session:
                query_data = session.query(WebLoginVerification.user_id).filter(WebLoginVerification.code == code).first()
            if query_data is None:
                return False, None
            else:
                user_id = query_data[0]
                if user_id is None:
                    return False, None
                else:
                    return True, user_id
        else:
            return False, None

    def delete_code(self, code: str):
        with self._session_scope() as session:
            fetch_num = session.query(WebLoginVerification).filter(
                WebLoginVerification.code == code).delete(synchronize_session=False)
            session.commit()
        return fetch_num

    def delete_expiration_code(self):
        with self._session_scope() as session:
            fetch_num = session.query(WebLoginVerification).filter(
                datetime.datetime.now() > WebLoginVerification.expiration_time).delete(synchronize_session=False)
            session.commit()
        return fetch_num

    def is_user_info_exist(self, user_id: str):
        with self._session_scope() as session:
            query_data = session.query(WebUserInfo.user_id).filter(WebUserInfo.user_id == user_id).first()
        if query_data is None:
            return False
        else:
            return True

    def add_user_info(self, user_id: str, name: str, avatar_url: str):
        with self._session_scope() as session:
            web_user_info = WebUserInfo(user_id, name, avatar_url)
            session.add(web_user_info)
            session.commit()

    def get_user_info(self, user_id):
        with self._session_scope() as session:
            query_data = session.query(WebUserInfo.name, WebUserInfo.avatar_url).filter(WebUserInfo.user_id == user_id).first()
        if query_data is None:
            return None
        else:
            name = query_data[0]
            avatar_url = query_data[1]
            return name, avatar_url


if __name__ == '__main__':
    pass
