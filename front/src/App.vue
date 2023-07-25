<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="menu-wrapper">
    <router-link class="item-link" activeClass="active-item-link" to="/">Monitor Interface</router-link>
    <router-link class="item-link" activeClass="active-item-link" to="/lan">Monitor Lan</router-link>
  </div>
  <RouterView />
</template>

<style lang="less" scoped>
.menu-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .item-link {
    display: block;
    font-size: 18px;
    font-weight: bold;
    color: #fff;
    text-decoration: none;
    text-align: center;
    flex: 1;
    height: 56px;
    line-height: 56px;
  }

  .active-item-link {
    color: #10b981
  }
}
</style>

