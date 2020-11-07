<template>
  <q-card class="my-card col-xs-6 col-sm-4 col-md-3 col-lg-2 shadow-1">
    <q-img :src="img_url" :ratio="16 / 9" spinner-color="orange">
      <div class="absolute-full flex flex-center text-h6" v-if="to_delete">
        Delete
      </div>
    </q-img>
    <q-card-section>
      <div class="row justify-between">
        <div class="text-h6" :class="{ 'text-strike': to_delete }">
          ID：{{ img_id }}
        </div>
        <div>
          <q-btn
            :flat="!is_gif"
            round
            :glossy="is_gif"
            color="light-blue-5"
            :text-color="gif_btn_color"
            :icon="mdiGif"
            :disable="to_delete"
            @click="gif_click"
          />
          <q-btn
            :flat="!to_delete"
            round
            :glossy="to_delete"
            color="light-blue-5"
            :text-color="del_btn_color"
            icon="delete"
            @click="del_click"
          />
        </div>
      </div>
    </q-card-section>
    <q-card-section class="row q-gutter-md justify-between">
      <q-input
        class="col-8"
        v-model.lazy="lineedit_url"
        placeholder="Image Url"
        :dense="lineedit_dense"
        :disable="to_delete"
      />
      <q-btn
        flat
        round
        class="col-2"
        :text-color="chech_btn_color"
        icon="check"
        :disable="to_delete"
        @click="loadImg(lineedit_url)"
      />
    </q-card-section>
  </q-card>
</template>

<script>
import { mdiGif } from "@quasar/extras/mdi-v4";
//import { matGif } from "@quasar/extras/material-icons";

var tri_color = "black";
var not_tri_color = "grey";

export default {
  name: "img-card",
  props: ["img_id", "input_url", "input_is_gif"],
  data: function() {
    return {
      img_url: this.input_url,
      lineedit_url: this.input_url,
      is_gif: this.input_is_gif,
      lineedit_dense: true,
      to_delete: false,
      gif_btn_color: "grey",
      del_btn_color: "grey",
      chech_btn_color: "green",
      mdiGif: mdiGif
    };
  },
  created: function() {
    /*
    this.img_url = this.input_url;
    this.lineedit_url = this.img_url;
    this.is_gif = this.input_is_gif;
    */
  },
  watch: {
    is_gif: function() {
      if (this.is_gif) {
        this.gif_btn_color = tri_color;
      } else {
        this.gif_btn_color = not_tri_color;
      }
      this.check_data_change();
    },
    to_delete: function() {
      if (this.to_delete) {
        this.del_btn_color = tri_color;
      } else {
        this.del_btn_color = not_tri_color;
      }
      this.check_data_change();
    },
    lineedit_url: function() {
      this.check_img();
      this.check_data_change();
    }
  },
  methods: {
    loadImg: function(input_url) {
      this.img_url = input_url;
      this.check_img();
      if (this.img_url.substring(this.img_url.length - 4) == ".gif") {
        this.is_gif = true;
      }
      this.check_data_change();
    },
    check_img() {
      if (this.lineedit_url == this.img_url) {
        this.chech_btn_color = "green";
      } else {
        this.chech_btn_color = not_tri_color;
      }
    },
    gif_click() {
      this.is_gif = !this.is_gif;
    },
    del_click: function() {
      this.to_delete = !this.to_delete;
    },
    check_data_change() {
      if (isNaN(this.img_id)) {
        if (this.to_delete) {
          this.emit_nochage_signal();
          return;
        }
        if (this.img_url == "") {
          this.emit_nochage_signal();
          return;
        } else {
          this.emit_chage_signal();
          return;
        }
      } else {
        if (this.to_delete) {
          this.emit_chage_signal();
          return;
        }
        if (this.is_gif != this.input_is_gif) {
          this.emit_chage_signal();
          return;
        }
        if (this.img_url != this.input_url) {
          this.emit_chage_signal();
          return;
        }

        this.emit_nochage_signal();
      }
    },
    emit_chage_signal() {
      this.$emit("icard-has-change", this.img_id);
    },
    emit_nochage_signal() {
      this.$emit("icard-no-change", this.img_id);
    }
  }
};
</script>

<style scoped>
.multi_bg_example {
  background-image: url(https://mdn.mozillademos.org/files/11307/bubbles.png),
    linear-gradient(to right, rgba(30, 75, 115, 1), rgba(0, 195, 255, 0.5));
  background-repeat: no-repeat, no-repeat;
  background-position: left, right;
}
</style>
