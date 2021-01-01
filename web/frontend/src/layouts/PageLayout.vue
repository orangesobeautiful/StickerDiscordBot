<template>
  <q-layout view="hHh lpR fFf">
    <q-header
      reveal
      :reveal-offset="100"
      elevated
      class="text-pink-4 multi_bg_example"
      height-hint="100"
    >
      <q-toolbar class="row justify-between">
        <q-toolbar-title class="row items-center">
          <div class="text-h6">
            (ノ・ω・)ノヾ(・ω・ヾ)
          </div>
        </q-toolbar-title>
        <q-avatar size="63px" class="q-ma-sm">
          <img :src="avatar_url" />
          <q-menu>
            <q-list style="min-width: 100px">
              <q-item clickable v-close-popup @click="logout">
                <q-item-section>登出</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-avatar>
      </q-toolbar>

      <q-tabs align="left">
        <q-route-tab to="/img-manage" label="貼圖管理" />
        <q-route-tab to="/sticker-web-tutorial" label="網頁版教學" />
        <q-route-tab
          to="/bot-instruction-tutorial"
          label="Bot指令教學(還沒有)"
        />
      </q-tabs>
    </q-header>

    <q-page-container>
      <q-page :style-fn="myTweak">
        <router-view />
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
export default {
  data() {
    return {
      name: "匿名",
      avatar_url: "../statics/404-avatar.png"
    };
  },
  async created() {
    var has_login = await this.check_has_login();
    if (!has_login) {
      // 如果沒登入則導向登入介面
      this.$router.push("/login");
    }
    this.getUserInfo();
  },
  methods: {
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
    myTweak(offset) {
      return {
        minWidth: "450px",
        minHeight: offset ? `calc(100vh - ${offset}px)` : "100vh"
      };
    },
    getUserInfo() {
      var path = "/sndata/user_info";
      this.$axios.get(path).then(res => {
        this.name = res.data.name;
        this.avatar_url = res.data.avatar_url;
      });
    },
    async logout() {
      var path = "/sndata/logout";
      await this.$axios
        .get(path)
        .then(() => {
          this.$router.push("/success-logout");
        })
        .catch(() => {});
    }
  }
};
</script>

<style scoped>
.multi_bg_example {
  background-image: url(../statics/snowdrop-1025050_1920_dark.jpg);
}

.he200 {
  height: 200px;
}
</style>
