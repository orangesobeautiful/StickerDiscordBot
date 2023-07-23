<template>
  <div class="column">
    <div class="row q-pa-lg justify-center q-gutter-md">
      <q-input
        square
        filled
        v-model="search_query"
        label="貼圖名稱"
        @keydown.enter.prevent="search"
        class="col-5"
      />
      <q-btn icon="search" glossy color="purple" @click="search" />
    </div>
    <div>
      <q-separator inset color="orange" />
    </div>
    <div class="row q-pa-lg justify-center">
      <q-pagination
        v-model="current_page"
        :max="max_page"
        :input="true"
        color="green"
        input-class="text-red"
        size="30px"
        @input="page_value_change"
      ></q-pagination>
    </div>
    <div class="row q-pl-xl q-pr-xl q-gutter-md">
      <sticker-group
        v-for="sn_group in sticker_list"
        :key="sn_group.key"
        :input_key="sn_group.key"
        :input_sticker_name="sn_group.sn"
        :img_list="sn_group.sts"
        :is_new_sg="sn_group.is_new_sg"
        @delete-sg="refresh"
        @add_success="to_normal_form"
      />
    </div>
    <div class="row q-pa-xl q-gutter-md">
      <q-btn
        class="self-start"
        color="amber"
        glossy
        label="New"
        @click="new_sticker"
      />
    </div>
    <div class="row q-gutter-md justify-center items-center">
      <q-pagination
        v-model="current_page"
        :max="max_page"
        :input="true"
        color="green"
        input-class="text-red"
        size="30px"
        @input="page_value_change"
      ></q-pagination>
    </div>
  </div>
</template>

<script>
import StickerGroup from "../components/StickerGroup.vue";

export default {
  name: "img-manage",
  components: {
    "sticker-group": StickerGroup
  },
  data: function() {
    return {
      sticker_list: Array(),
      current_page: 1,
      latest_get_page: 0,
      page_show_num: 10,
      max_page: 1,
      new_sg_num: 0,
      search_query: ""
    };
  },
  created: function() {
    this.get_stickers();
  },
  methods: {
    page_value_change() {
      if (this.current_page != this.latest_get_page) {
        this.get_stickers();
      }
    },
    get_stickers() {
      // eslint-disable-next-line
      //console.log(this.get_stickers);
      //const path = "http://localhost:" + this.GLOBAL.BACKENDPORT.toString() + "/all_sticker";
      var path = "/all_sticker";
      this.$axios
        .get(path, {
          params: {
            start: (this.current_page - 1) * this.page_show_num,
            num: this.page_show_num
          }
        })
        .then(res => {
          var img_data = res.data.img_data;
          for (var i = 0; i < img_data.length; i++) {
            img_data[i].sts.sort(function(a, b) {
              return a.id - b.id;
            });
            img_data[i].key = img_data[i].sn;
          }
          this.sticker_list = img_data;
          this.max_page = res.data.maxp;
          this.latest_get_page = this.current_page;
        });
    },
    refresh() {
      this.get_stickers();
    },
    new_sticker() {
      this.sticker_list.push({
        sn: "",
        sts: Array(),
        is_new_sg: true,
        key: "NewSG" + this.new_sg_num.toString()
      });
      this.new_sg_num++;
    },
    to_normal_form(sg_key, sn) {
      var sg_index = this.sticker_list
        .map(function(item) {
          return item.key;
        })
        .indexOf(sg_key);
      this.sticker_list[sg_index].is_new_sg = false;
      this.sticker_list[sg_index].key = sn;
      this.sticker_list[sg_index].sn = sn;
    },
    search() {
      if (this.search_query == "") {
        this.get_stickers();
      } else {
        var path = "/search";
        this.$axios
          .get(path, {
            params: {
              q: this.search_query
            }
          })
          .then(res => {
            var img_data = res.data.img_data;
            if (img_data == null) {
              img_data = [];
            }
            for (var i = 0; i < img_data.length; i++) {
              img_data[i].sts.sort(function(a, b) {
                return a.id - b.id;
              });
              img_data[i].key = img_data[i].sn;
            }
            this.sticker_list = img_data;
            this.max_page = res.data.maxp;
            this.latest_get_page = this.current_page;
          });
      }
    }
  }
};
</script>
