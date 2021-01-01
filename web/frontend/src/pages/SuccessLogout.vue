<template>
  <div class="col justify-center">
    <h3 class="row justify-center">
      成功登出
    </h3>
    <a class="row justify-center" href="/login">返回登入介面</a>
  </div>
</template>

<script>
export default {
  name: "login",
  components: {},
  data() {
    return {};
  },
  async created() {
    var has_login = await this.check_has_login();
    if (has_login) {
      // 如果已登入直接導向首頁
      this.$router.push("/");
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
    }
  }
};
</script>
