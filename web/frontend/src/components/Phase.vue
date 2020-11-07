<template>
  <div>
    <div class="row items-center q-mt-md q-ml-xl q-mr-xl">
      <q-icon name="star" size="sm" color="green-12 " />
      <h4 class="text-green-13">{{ title }}</h4>
    </div>
    <div class="q-mt-md q-ml-xl q-mr-xl context">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "phase",
  props: ["title"],
  data: function() {
    return {};
  }
};
</script>

<style scoped>
.context {
  font-family: sans-serif, arial;
  line-height: 1.6;
}
</style>
