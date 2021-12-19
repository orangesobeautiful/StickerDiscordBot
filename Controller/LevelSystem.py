import math
from Database import LevelSystemOperation
from sqlalchemy.engine.base import Engine
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
import datetime
import pytz

_const_lv_base_name = "lv_base"
_const_lv_exponent_coefficient_name = "lv_exponent_coefficient"
_const_lv_coefficient_name = "lv_coefficient"
_const_max_lv_name = "max_lv"

_const_normal_online_exp_reward_name = "normal_online_exp_reward"
_const_special_online_exp_reward_name = "special_online_exp_reward"
_const_weekday_special_time_start_name = "weekday_special_time_start"
_const_weekday_special_time_end_name = "weekday_special_time_end"
# signin mission
_const_daily_signin_mission_target_name = "daily_signin_mission_reward"
_const_daily_signin_mission_reward_name = "daily_signin_mission_reward"
# message mission
_const_daily_message_mission_target_name = "daily_message_mission_target"
_const_daily_message_mission_reward_name = "daily_message_mission_reward"
# word mission
_const_daily_word_mission_target_name = "daily_word_mission_target"
_const_daily_word_mission_reward_name = "daily_word_mission_reward"
# weekly mission setting
_const_weekly_mission_target_times_name = "weekly_mission_target_times"
_const_weekly_mission_reward_times_name = "weekly_mission_reward_times"


class LevelSystemController:
    def __init__(self, engine: Engine, timezone: pytz.UTC):
        self._session_maker = sessionmaker(bind=engine)
        self._timezone: pytz.UTC = timezone
        LevelSystemOperation.init_database_tables(engine)

        # read setting
        with self._session_maker() as session:
            lv_base: float = \
                float(LevelSystemOperation.init_lv_setting(session, _const_lv_base_name, str(1.05)))
            lv_coefficient: float = \
                float(LevelSystemOperation.init_lv_setting(session, _const_lv_coefficient_name, str(10000)))
            self.max_lv: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_max_lv_name, str(100)))
            self.lv_cumulative_exp_list: list = self.lv_exp_calculate(lv_base, lv_coefficient)

            self.normal_online_exp_reward: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_normal_online_exp_reward_name, str(1)))
            self.special_online_exp_reward: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_special_online_exp_reward_name, str(1 * 10)))
            self.weekday_special_time_start: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_weekday_special_time_start_name, str(17)))
            self.weekday_special_time_end: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_weekday_special_time_end_name, str(1)))

            self.daily_signin_mission_reward: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_daily_signin_mission_reward_name,
                                                         str(1 * 60 * 10 * 3)))
            self.daily_message_mission_target: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_daily_message_mission_target_name, str(5)))
            self.daily_message_mission_reward: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_daily_message_mission_reward_name,
                                                         str(1 * 60 * 10 * 6)))
            self.daily_word_mission_target: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_daily_word_mission_target_name, str(50)))
            self.daily_word_mission_reward: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_daily_word_mission_reward_name,
                                                         str(1 * 60 * 10 * 6)))

            self.weekly_mission_target_times: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_weekly_mission_target_times_name, str(4)))
            self.weekly_mission_reward_times: int = \
                int(LevelSystemOperation.init_lv_setting(session, _const_weekly_mission_reward_times_name, str(7)))

    def current_lv(self, current_cumulative_exp: int) -> int:
        """根據黨前累計經驗計算現在等級"""
        # 有時間用二分搜尋法
        for lv, cumulative_exp in enumerate(self.lv_cumulative_exp_list):
            if current_cumulative_exp < cumulative_exp:
                return lv
        return self.max_lv

    def lv_exp_calculate(self, lv_base: float,  lv_coefficient: float) -> list():
        """計算每個等級的累積經驗"""
        lv_exp_list = [0]
        lv_cumulative_exp_list = list()
        for lv in range(self.max_lv):
            lv_exp_list.append(lv_coefficient*math.pow(lv_base, lv))
        cumulative_exp = 0
        for lv_exp in lv_exp_list:
            cumulative_exp += lv_exp
            lv_cumulative_exp_list.append(int(cumulative_exp))

        return lv_cumulative_exp_list

    def is_special_time(self) -> bool:
        """判斷是否為掛機經驗特別時間"""
        now = datetime.datetime.now(self._timezone)
        if now.weekday() > 4:
            # 假日全天皆為 special time
            return True
        else:
            if self.weekday_special_time_end > self.weekday_special_time_start:
                if self.weekday_special_time_end > now.hour > self.weekday_special_time_start:
                    return True
            else:
                # 跨天需特別處理
                if now.hour > self.weekday_special_time_start or self.weekday_special_time_end > now.hour:
                    return True
        return False

    def get_user_exp_info(self, user_id: str, guild_id: str) -> LevelSystemOperation.UserEXP:
        """"get user exp info"""
        with self._session_maker() as session:
            user_exp_info = LevelSystemOperation.get_user_exp_info(session, user_id, guild_id)
        return user_exp_info

    def add_user_exp(self, user_id: str, guild_id: str, exp: int):
        """add user exp"""
        with self._session_maker() as session:
            LevelSystemOperation.add_exp(session, user_id, guild_id, exp)

    def emit_online_exp_reward(self, guild_online_members: dict) -> None:
        """發送掛機經驗值"""
        exp_reward = self.normal_online_exp_reward
        if self.is_special_time():
            exp_reward = self.special_online_exp_reward

        for guild_id in guild_online_members.keys():
            user_list = guild_online_members[guild_id]
            for user_id in user_list:
                self.add_user_exp(user_id, guild_id, exp_reward)

    def daily_reset(self) -> None:
        """重製每日和每周任務狀態"""
        now = datetime.datetime.now(self._timezone)
        with self._session_maker() as session:
            reset_week = False
            if now.weekday() == 0:
                reset_week = True
            LevelSystemOperation.reset_mission_status(session, reset_week)

    def user_signin(self, user_id: str, guild_id: str) -> str:
        """訊息簽到狀態"""
        user_exp_info = self.get_user_exp_info(user_id, guild_id)
        update_col = dict()
        mission_reward = 0
        exp_status_msg = ""

        # 檢查簽到狀態
        if not user_exp_info.daily_signin:
            # 更新每日簽到並發放經驗值
            update_col[LevelSystemOperation.UserEXP.daily_signin] = True
            emit_exp = self.daily_signin_mission_reward
            mission_reward += emit_exp
            exp_status_msg = "每日簽到任務完成，獲得 " + str(emit_exp) + " 經驗值\n"
            # 檢查每周簽到任務
            if user_exp_info.weekly_signin < 1 * self.weekly_mission_target_times:
                user_exp_info.weekly_signin += 1
                update_col[LevelSystemOperation.UserEXP.weekly_signin] = user_exp_info.weekly_signin
                # 每周簽到任務達標，發放經驗值
                if user_exp_info.weekly_signin == 1 * self.weekly_mission_target_times:
                    emit_exp = self.daily_signin_mission_reward * self.weekly_mission_reward_times
                    mission_reward += emit_exp
                    exp_status_msg = "每周簽到任務完成，獲得 " + str(emit_exp) + " 經驗值\n"

        # 更新經驗值
        if mission_reward > 0:
            update_col[LevelSystemOperation.UserEXP.exp] = LevelSystemOperation.UserEXP.exp + mission_reward

        # 更新資料
        if len(update_col) > 0:
            with self._session_maker() as session:
                LevelSystemOperation.update_user_exp_info(session, user_id, guild_id, update_col)

        return exp_status_msg

    def message_deal(self, user_id: str, guild_id: str, msg_content: str, is_command: bool) -> str:
        """訊息處理，當 member 發出訊息透過此函示處理(每日任務、每周任務)"""
        user_exp_info = self.get_user_exp_info(user_id, guild_id)
        update_col = dict()
        mission_reward = 0
        exp_status_msg = ""

        exp_status_msg += self.user_signin(user_id, guild_id)
        if not is_command:
            # 檢查每日訊息量任務
            if user_exp_info.daily_message < self.daily_message_mission_target:
                user_exp_info.daily_message += 1
                update_col[LevelSystemOperation.UserEXP.daily_message] = user_exp_info.daily_message
                # 每日訊息量任務達標，發放經驗值
                if user_exp_info.daily_message == self.daily_message_mission_target:
                    emit_exp = self.daily_message_mission_reward
                    mission_reward += emit_exp
                    exp_status_msg = "每日訊息量任務完成，獲得 " + str(emit_exp) + " 經驗值\n"
            # 檢查每日字數量任務
            if user_exp_info.daily_word < self.daily_word_mission_target:
                user_exp_info.daily_word += len(msg_content)
                if user_exp_info.daily_word > self.daily_word_mission_target:
                    update_col[LevelSystemOperation.UserEXP.daily_word] = self.daily_word_mission_target
                else:
                    update_col[LevelSystemOperation.UserEXP.daily_word] = user_exp_info.daily_word
                # 每日字數量任務達標，發放經驗值
                if user_exp_info.daily_word == self.daily_word_mission_target:
                    emit_exp = self.daily_word_mission_reward
                    mission_reward += emit_exp
                    exp_status_msg = "每日字數量任務完成，獲得 " + str(emit_exp) + " 經驗值\n"
            # 檢查每週訊息量任務
            if user_exp_info.weekly_message < self.daily_message_mission_target * self.weekly_mission_target_times:
                user_exp_info.weekly_message += 1
                update_col[LevelSystemOperation.UserEXP.weekly_message] = user_exp_info.weekly_message
                # 每週訊息量任務達標，發放經驗值
                if user_exp_info.weekly_message == self.daily_message_mission_target * self.weekly_mission_target_times:
                    emit_exp = self.daily_message_mission_reward * self.weekly_mission_reward_times
                    mission_reward += emit_exp
                    exp_status_msg = "每周訊息量任務完成，獲得 " + str(emit_exp) + " 經驗值\n"
            # 檢查每週字數量任務
            if user_exp_info.weekly_word < self.daily_word_mission_target * self.weekly_mission_target_times:
                user_exp_info.weekly_word += len(msg_content)
                if user_exp_info.weekly_word > self.daily_word_mission_target * self.weekly_mission_target_times:
                    update_col[LevelSystemOperation.UserEXP.weekly_word] = self.daily_word_mission_target * self.weekly_mission_target_times
                else:
                    update_col[LevelSystemOperation.UserEXP.weekly_word] = user_exp_info.weekly_word
                # 每週字數量任務達標，發放經驗值
                if user_exp_info.weekly_word == self.daily_word_mission_target * self.weekly_mission_target_times:
                    emit_exp = self.deaily_word_mission_reward * self.weekly_mission_reward_times
                    mission_reward += emit_exp
                    exp_status_msg = "每周字數量任務完成，獲得 " + str(emit_exp) + " 經驗值\n"

        # 更新經驗值
        if mission_reward > 0:
            update_col[LevelSystemOperation.UserEXP.exp] = LevelSystemOperation.UserEXP.exp + mission_reward

        # 更新資料
        if len(update_col) > 0:
            with self._session_maker() as session:
                LevelSystemOperation.update_user_exp_info(session, user_id, guild_id, update_col)

        return exp_status_msg

    def user_lv_info(self, user_name: str, user_id: str, guild_id: str) -> str:
        exp_info: LevelSystemOperation.UserEXP = self.get_user_exp_info(user_id, guild_id)
        if exp_info is None:
            return user_name + "目前還沒有等級資訊"
        if exp_info.daily_signin:
            signin_show_str = "✓"
        else:
            signin_show_str = "✗"

        next_lv_deviation = 0
        current_lv = self.current_lv(exp_info.exp)
        if current_lv != self.max_lv:
            next_lv_deviation = self.lv_cumulative_exp_list[current_lv] - exp_info.exp

        return_msg = ""
        return_msg += user_name + " 的等級資訊:\n\n"
        return_msg += f"LV:\t{current_lv}\n"
        return_msg += f"EXP:\t{exp_info.exp}\n"
        return_msg += f"距離下個等級:\t{next_lv_deviation}\n"
        return_msg += "\n"
        return_msg += f"每日簽到任務:\t{signin_show_str}\n"
        return_msg += f"每日訊息量任務:\t{exp_info.daily_message}/{self.daily_message_mission_target}\n"
        return_msg += f"每日字數量任務:\t{exp_info.daily_word}/{self.daily_word_mission_target}\n"
        return_msg += f"每周簽到任務:\t{exp_info.weekly_signin}/{self.weekly_mission_target_times}\n"
        return_msg += f"每周訊息量任務:\t{exp_info.weekly_message}/{self.daily_message_mission_target*self.weekly_mission_target_times}\n"
        return_msg += f"每周字數量任務:\t{exp_info.weekly_word}/{self.daily_word_mission_target*self.weekly_mission_target_times}\n"
        return return_msg


if __name__ == '__main__':
    _test_db_url = 'mysql+pymysql://test:1234@localhost/our_bot'
    _test_engine = create_engine(
        _test_db_url, pool_pre_ping=True, echo=False, pool_recycle=7200)
    _test_c = LevelSystemController(_test_engine, pytz.timezone('Asia/Taipei'))
    _test_c.is_special_time()


