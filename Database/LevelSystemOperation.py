import sqlalchemy
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from sqlalchemy.orm.session import Session
from sqlalchemy.engine.base import Engine
import time

_Base = declarative_base()

_lv_setting_table_name = 'lv_setting'
_lv_name_table_name = 'lv_name'
_user_exp_table_name = 'user_exp'


class LVName(_Base):
    __tablename__ = _lv_name_table_name

    lv = sqlalchemy.Column('lv', sqlalchemy.Integer, primary_key=True, nullable=False)
    name = sqlalchemy.Column('name', sqlalchemy.Text, nullable=False)

    def __init__(self, lv: int, name: str):
        self.lv = lv
        self.name = name

    def __repr__(self):
        return "<LVName('{}', '{}')>".format(self.lv, self.name)

    def items(self):
        return [self.lv, self.name]

    def create(self, session):
        session.add(self)
        session.commit()


class LVSetting(_Base):
    __tablename__ = _lv_setting_table_name

    name = sqlalchemy.Column('name', sqlalchemy.String(
        128), primary_key=True, nullable=False)
    value = sqlalchemy.Column('value', sqlalchemy.Text, nullable=False)

    def __init__(self, name: str, value: str):
        self.name = name
        self.value = value

    def __repr__(self):
        return "<LVSetting('{}', '{}')>".format(self.name, self.value)

    def items(self):
        return [self.name, self.value]

    def create(self, session):
        session.add(self)
        session.commit()


class UserEXP(_Base):
    __tablename__ = _user_exp_table_name

    user_id = sqlalchemy.Column(
        'user_id', sqlalchemy.String(128), primary_key=True, nullable=False)
    guild_id = sqlalchemy.Column(
        'guild_id', sqlalchemy.String(128), primary_key=True ,nullable=False)
    exp = sqlalchemy.Column(
        'exp', sqlalchemy.BigInteger, default=0, nullable=False)
    daily_signin = sqlalchemy.Column(
        'daily_signin', sqlalchemy.Boolean, default=False, nullable=False)
    daily_message = sqlalchemy.Column(
        'daily_message', sqlalchemy.Integer, default=0, nullable=False)
    daily_word = sqlalchemy.Column(
        'daily_word', sqlalchemy.Integer, default=0, nullable=False)
    weekly_signin = sqlalchemy.Column(
        'weekly_signin', sqlalchemy.Integer, default=0, nullable=False)
    weekly_message = sqlalchemy.Column(
        'weekly_message', sqlalchemy.Integer, default=0, nullable=False)
    weekly_word = sqlalchemy.Column(
        'weekly_word', sqlalchemy.Integer, default=0, nullable=False)
    create_at = sqlalchemy.Column(
        'create_at', sqlalchemy.BigInteger, default=int(time.time()), nullable=False)

    def __init__(self, user_id: str, guild_id: str):
        self.user_id = user_id
        self.guild_id = guild_id

    def __repr__(self):
        return "<UserEXP('{}', '{}', '{}')>".format(self.user_id, self.guild_id, self.exp)

    def items(self):
        return [self.user_id, self.guild_id, self.exp]

    def create(self, session: Session):
        session.add(self)
        session.commit()


def init_database_tables(engine: Engine):
    """check table is exist
    if not exist than create """
    insp = sqlalchemy.inspect(engine)
    # if table is not exist than create
    if not insp.has_table(_user_exp_table_name, None):
        UserEXP.metadata.create_all(engine)

    if not insp.has_table(_lv_setting_table_name, None):
        LVSetting.metadata.create_all(engine)

    if not insp.has_table(_lv_name_table_name, None):
        LVName.metadata.create_all(engine)


def get_lv_setting(session: Session, name: str):
    query_data = session.query(LVSetting.value).filter(
        LVSetting.name == name).first()
    if query_data is None:
        return None
    return query_data[0]


def set_lv_setting(session: Session, name: str, value: str):
    """set lv setting"""
    session.query(LVSetting).filter(
        LVSetting.name == name).update({LVSetting.value: value})
    session.commit()


def upsert_lv_setting(session: Session, name: str, value: str):
    """update lv setting, if lv setting is not exit, than create one"""
    if get_lv_setting(session, name) is None:
        LVSetting(name, value).create(session)
    else:
        set_lv_setting(session, name, value)


def init_lv_setting(session: Session, name: str, default_value: str) -> str:
    """read lv setting, if setting is not exist, than create by default value"""
    value = get_lv_setting(session, name)
    if value is None:
        value = default_value
        LVSetting(name, default_value).create(session)
    return value


def get_user_exp_info(session: Session, user_id: str, guild_id: str) -> UserEXP:
    """get user exp info"""
    exp_info = session.query(UserEXP).filter(
        sqlalchemy.and_(UserEXP.user_id == user_id, UserEXP.guild_id == guild_id)).first()
    if exp_info is None:
        user_exp_info = UserEXP(user_id, guild_id)
        user_exp_info.create(session)
        return user_exp_info
    return exp_info


def update_user_exp_info(session: Session, user_id: str, guild_id: str, update_col: dict) -> None:
    """update user exp info by given value"""
    session.query(UserEXP).filter(
        sqlalchemy.and_(UserEXP.user_id == user_id, UserEXP.guild_id == guild_id)).update(update_col)
    session.commit()


def add_exp(session: Session, user_id: str, guild_id: str, exp: int) -> None:
    """add user exp, is user is not exist will create"""
    session.query(UserEXP).filter(sqlalchemy.and_(
        UserEXP.user_id == user_id, UserEXP.guild_id == guild_id)).update(
        {UserEXP.exp: UserEXP.exp + exp})
    session.commit()


def reset_mission_status(session: Session, reset_week: bool) -> None:
    update_col = {
        UserEXP.daily_signin: False, UserEXP.daily_word: 0, UserEXP.daily_message: 0}
    if reset_week:
        update_col.update({
            UserEXP.weekly_word: 0, UserEXP.weekly_message: 0, UserEXP.weekly_signin: 0})

    session.query(UserEXP).update(update_col)
    session.commit()


if __name__ == '__main__':
    _test_db_url = 'mysql+pymysql://test:1234@localhost/our_bot'
    _test_engine = create_engine(
        _test_db_url, pool_pre_ping=True, echo=False, pool_recycle=7200)

    _test_session_maker = sessionmaker(bind=_test_engine)
    _test_session = _test_session_maker()
    set_lv_setting(_test_session, "test-name", "test-value")
    print(get_user_exp_info(_test_session, '123', '456'))
    print(get_lv_setting(_test_session, "test-name"))
    print(get_user_exp_info(_test_session, '23130', 'fdsfds'))
