<template>
  <div class="col">
    <div class="row justify-center">
      <h5>
        請在可認證的群組使用指令
      </h5>
      <h5 class="text-deep-orange-5 ">&nbsp;web-login &lt;驗證碼&gt;&nbsp;</h5>
      <h5>進行身分驗證</h5>
    </div>
    <div class="row justify-center">
      <div class="text-h3 q-pa-md" v-if="code_effective_time > 0">
        {{ login_code }}
      </div>
      <q-btn
        class="q-pa-md text-weight-bolder"
        color="light-green-14 "
        outline
        icon="refresh"
        v-if="code_effective_time <= 0"
        label="重新產生驗證碼"
        @click="get_valid_code"
      />
      <div
        class="detailed-color column flex-center bg-grey-5  q-ma-md q-pl-md q-pr-md non-selectable"
      >
        {{ code_effective_time_show_str }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "login",
  components: {},
  data() {
    return {
      login_code: -1,
      check_login_timer: null,
      code_effective_timer: null,
      code_effective_time: 0,
      code_effective_time_show_str: "0:00"
    };
  },
  async created() {
    var has_login = await this.check_has_login();
    if (has_login) {
      // 如果已登入直接導向首頁
      this.$router.push("/");
    } else {
      clearInterval(this.check_login_timer);
      this.check_login_timer = null;
      clearInterval(this.code_effective_timer);
      this.code_effective_timer = null;
      this.get_valid_code();
    }
  },
  methods: {
    // 檢查是否登入過
    async check_has_login() {
      var path = "/sndata/has_login";
      var has_login;
      await this.$axios
        .get(path)
        .then(() => {
          has_login = true;
        })
        .catch(() => {
          has_login = false;
        });
      return has_login;
    },
    // 獲取登入用認證碼
    get_valid_code() {
      var path = "/sndata/get_login_code";
      this.$axios.get(path).then(res => {
        this.login_code = res.data.code;
      });
      this.code_effective_time = 300;
      this.set_code_effective_timer();
      this.set_check_login_timer();
    },
    // 檢查認證碼登入狀態
    chech_login() {
      var path = "/sndata/check_login";
      this.$axios.get(path, { params: { code: this.login_code } }).then(res => {
        if (res.data.result == "1") {
          this.$router.push("/");
        }
      });
    },
    // 設置 chech_login 計時器
    set_check_login_timer() {
      if (this.check_login_timer == null) {
        this.check_login_timer = setInterval(() => {
          this.chech_login();
        }, 3000);
      }
    },
    set_code_effective_timer() {
      if (this.code_effective_timer == null) {
        this.code_effective_timer = setInterval(() => {
          this.code_effective_time -= 1;
          if (this.code_effective_time <= 0) {
            this.valid_code_expired();
          } else {
            var min = parseInt(this.code_effective_time / 60);
            var sec = this.code_effective_time % 60;
            var sec_zero_padding = "";
            if (sec < 10) {
              sec_zero_padding = "0";
            }
            sec_zero_padding += sec.toString();

            this.code_effective_time_show_str = min + ":" + sec_zero_padding;
          }
        }, 1000);
      }
    },
    valid_code_expired() {
      clearInterval(this.code_effective_timer);
      this.code_effective_timer = null;
      this.code_effective_time = 0;
    }
  },
  destroyed: function() {
    clearInterval(this.check_login_timer);
    this.check_login_timer = null;
    clearInterval(this.code_effective_timer);
    this.code_effective_timer = null;
  }
};
</script>
