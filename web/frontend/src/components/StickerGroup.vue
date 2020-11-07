<template>
  <q-card bordered class="shadow-12 col-12 whole_bg" style="">
    <q-card-section class="text-black col-12 row  justify-center">
      <div v-if="!is_new_sg" class="text-h5">{{ sticker_name }}</div>
      <q-input v-if="is_new_sg" v-model="sticker_name" label="貼圖" />
    </q-card-section>

    <q-separator spaced inset />

    <q-card-section class="row justify-center q-pa-lg q-gutter-xl" style="">
      <img-card
        v-for="img in img_list"
        :input_url="img.url"
        :img_id="img.id"
        :input_is_gif="img.gif"
        :key="'img' + img.id"
        :ref="'img' + img.id"
        @icard-has-change="new_data_change"
        @icard-no-change="un_data_change"
      ></img-card>
    </q-card-section>
    <q-card-section class="row justify-between">
      <q-btn round icon="add" color="green-6" size="lg" @click="add_img" />
      <q-btn
        push
        label="OK"
        color="blue-6"
        size="lg"
        :disable="!enable_ok_btn"
        @click="apply_change"
      />
    </q-card-section>
  </q-card>
</template>

<script>
import ImgCard from "./ImgCard.vue";

export default {
  name: "sticker-group",
  props: ["input_key", "input_sticker_name", "img_list", "is_new_sg"],
  components: {
    "img-card": ImgCard
  },
  data: function() {
    return {
      new_img_num: 1,
      sticker_name: this.input_sticker_name,
      change_id_list: Array(),
      all_change: {},
      enable_ok_btn: false
    };
  },
  created: function() {
    if (this.is_new_sg) this.add_img();
  },
  methods: {
    add_img() {
      this.img_list.push({
        id: "New Image" + this.new_img_num.toString(),
        url: "",
        gif: false
      });
      this.new_img_num++;
    },
    new_data_change(id) {
      this.change_id_list.push(id);
      this.enable_ok_btn = true;
    },
    un_data_change(id) {
      this.change_id_list = this.change_id_list.filter(e => e !== id);
      if (this.change_id_list.length == 0) {
        this.enable_ok_btn = false;
      }
    },
    apply_change() {
      var has_change = true;
      this.all_change = {
        sn: this.sticker_name,
        add: Array(),
        edit: Array(),
        delete: Array()
      };
      this.img_list.forEach(this.detect_change);
      if (has_change) {
        var path = "/sndata/change_sn";
        this.$axios({
          method: "post",
          url: path,
          data: this.all_change
        })
          .then(res => {
            if (this.is_new_sg) {
              this.$emit("add_success", this.input_key, this.sticker_name);
            }
            this.img_list.length = 0;
            var new_img_data = res.data["imgs"];

            if (new_img_data.length == 0) {
              this.$emit("delete-sg");
            }

            for (var i = 0; i < new_img_data.length; i++) {
              this.img_list.push({
                id: new_img_data[i][0],
                url: new_img_data[i][1],
                gif: new_img_data[i][2]
              });
            }
            this.change_id_list.length = 0;
            this.enable_ok_btn = false;
          })
          .catch(error => {
            // eslint-disable-next-line
            console.log(error);
          });
      }
    },
    detect_change(img_item) {
      var img_ref = this.$refs["img" + img_item.id][0];
      var img_id = img_item.id;
      var input_url = img_item.url;
      var to_delete = img_ref.to_delete;
      var img_url = img_ref.img_url;
      var input_is_gif = img_item.gif;
      var is_gif = img_ref.is_gif;

      // eslint-disable-next-line
      //console.log(img_ref);

      // eslint-disable-next-line

      if (isNaN(img_id)) {
        // eslint-disable-next-line
        if (!to_delete)
          this.all_change["add"].push({ url: img_url, gif: is_gif });
      } else if (to_delete) {
        this.all_change["delete"].push(parseInt(img_id));
      } else {
        if (img_url != input_url) {
          this.all_change["edit"].push({ id: img_id, url: img_url });
        }
        if (is_gif != input_is_gif) {
          this.all_change["edit"].push({ id: img_id, gif: is_gif });
        }
      }
    }
  }
};
</script>

<style scoped>
.card {
  background: radial-gradient(circle, #35a2ff 20%, #014a88 100%);
}

.whole_bg {
  background-image: url("../statics/bird-388258.svg"),
    linear-gradient(to right, rgb(228, 81, 81), rgba(255, 145, 0, 0.5));
  background-size: 150px, cover;
  background-repeat: no-repeat, no-repeat;
  background-position: 100% 30%, right;
}
</style>
