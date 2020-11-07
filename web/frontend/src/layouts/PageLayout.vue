<template>
  <q-layout view="hHh lpR fFf">
    <q-header
      reveal
      :reveal-offset="100"
      elevated
      class="text-pink-4 multi_bg_example"
      height-hint="100"
    >
      <q-toolbar>
        <q-toolbar-title class="row items-center">
          <q-avatar size="63px" class="q-ma-sm">
            <img src="../statics/1b669b15ea55d1530a2c352d15f4b7fa.png" />
          </q-avatar>
          <div class="text-h6">
            (ノ・ω・)ノヾ(・ω・ヾ)
          </div>
        </q-toolbar-title>
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
    return {};
  },
  methods: {
    myTweak(offset) {
      return {
        minWidth: "450px",
        minHeight: offset ? `calc(100vh - ${offset}px)` : "100vh"
      };
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
