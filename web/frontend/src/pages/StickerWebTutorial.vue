<template>
  <div>
    <div class="row justify-center q-gutter-m">
      <h3 class="text-deep-orange-5">網頁版貼圖管理系統教學</h3>
    </div>
    <phase title="介面">
      <q-img
        src="../statics/StickerWebTutorial/StickerGorup.png"
        style="width: 100%;"
        spinner-color="orange"
      >
      </q-img>
      <div class="text-h6">
        1.貼圖名稱<br />
        2.圖片<br />
        3.貼圖ID:同個名稱的貼圖可能有好幾張圖片，ID用來區分貼圖(不重要)<br />
        4.刪除按鈕(可勾選)<br />
        5.GIF按鈕(可勾選)<br />
        6.貼圖網址<br />
        7.變更圖片按鈕:貼圖網址改變後可使用，按下後圖片會預覽出當前網址的圖片<br />
        8.新增貼圖按鈕:是新增「同名稱」的貼圖圖片<br />
        9.確認變更按鈕:要新增、修改、刪除都要按下這個按鈕才會真的改變<br />
      </div>
    </phase>
    <phase title="新增貼圖">
      <h5>a.新名稱貼圖(新增完整的貼圖)</h5>
      <q-img
        src="../statics/StickerWebTutorial/AddNewNameSticker.png"
        style="width: 100%;"
        spinner-color="orange"
      />
      <div class="text-h6">
        1.按下網頁最左下的"New"按鈕，此時會有個空白貼圖群組<br />
        2.填入貼圖名稱(可有空白字符)<br />
        3.填入貼圖的圖片網址<br />
        4.按下勾勾按鈕，此時會預覽出圖片的樣子<br />
        5.確認沒問題後按下OK按鈕就新增完成啦<br />
      </div>
      <h5>b.新圖片貼圖(現有貼圖中新增圖片)</h5>
      <q-img
        src="../statics/StickerWebTutorial/AddNewImgSticker.png"
        style="width: 100%;"
        spinner-color="orange"
      />
      <div class="text-h6">
        1.按下左下的+形狀按鈕，此時會有個空白貼圖<br />
        2.填入貼圖的圖片網址<br />
        3.按下勾勾按鈕，此時會預覽出圖片的樣子<br />
        5.確認沒問題後按下OK按鈕就新增完成啦<br />
      </div>
    </phase>
    <phase title="修改貼圖">
      <h5>a.修改圖片</h5>
      <q-img
        src="../statics/StickerWebTutorial/EditImgSticker.png"
        style="width: 100%;"
        spinner-color="orange"
      />
      <div class="text-h6">
        1.變更想要修改貼圖的網址<br />
        2.按下勾勾按鈕<br />
        3.確認沒問題後按下OK按鈕就修改成功啦<br />
      </div>
    </phase>
    <phase title="刪除貼圖">
      <q-img
        src="../statics/StickerWebTutorial/DeleteImgSticker.png"
        style="width: 100%;"
        spinner-color="orange"
      />
      <div class="text-h6">
        1.按下想要刪除圖片的刪除按鈕<br />
        2.確認沒選作後按下OK按鈕就刪除啦<br />
        <div class="text-red">ps:把所有圖片刪除後即可刪除整個貼圖</div>
        <br />
      </div>
    </phase>
    <div class="row justify-center  q-mt-md q-ml-xl q-mr-xl q-mb-md"></div>
  </div>
</template>

<script>
import Phase from "../components/Phase.vue";

export default {
  name: "sticker-web-tutorial",
  components: { phase: Phase },
  data: function() {
    return {};
  }
};
</script>
