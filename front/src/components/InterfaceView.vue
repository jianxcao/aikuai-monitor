<script setup lang="ts">
import { FormatSize, FormatNum } from '@/utils'
import { computed } from 'vue';
const props = defineProps({
  data: Object
})
const linkData = computed(() => {
  const data = props.data || {}
  return Object.keys(data).reduce((prev: Record<string, any>, key) => {
    const d = data[key].data.ifaceStream
    prev[key] = d.reduce((res: Record<string, number>, item: any) => {
      const connectNum = Number(item.connectNum)
      if (connectNum) {
        res.connectNum += connectNum
        res.download += item.download
        res.upload += item.upload
      }
      return res;
    }, {
      download: 0,
      upload: 0,
      connectNum: 0
    })
    return prev;
  }, {})
})
</script>

<template>
  <div class="interfaceWrapper" v-for="(item, key) in data" :key="key">
    <div class="header">{{ key }}</div>
    <div class="interface">
      <div class="left">
        <div class="speed">
          <div class="title">当前上传</div>
          <div class="content"><span class="uploadTxt">{{ FormatSize(linkData[key]?.upload) }}/S</span></div>
        </div>
        <div class="speed">
          <div class="title">当前下载</div>
          <div class="content"><span class="downloadTxt">{{ FormatSize(linkData[key]?.download) }}/S</span></div>
        </div>
        <div class="speed">
          <div class="title">总连接数</div>
          <div class="content"><span class="downloadTxt">{{ FormatNum(linkData[key]?.connectNum) }}</span></div>
        </div>
      </div>
      <div class="interfaceGirdData right">
        <table class="ifaceCheck">
          <colgroup>
            <col>
            <col>
            <col>
            <col>
            <col>
          </colgroup>
          <thead>
            <tr>
              <th>连接名称</th>
              <th>IP</th>
              <th>联网方式</th>
              <th>结果</th>
              <th>备注</th>
            </tr>
          </thead>
          <tbody v-if="item?.data?.ifaceCheck && item.data.ifaceCheck.length">
            <tr v-for="(ifaceItem, index) in item.data.ifaceCheck" :key="index">
              <td>
                {{ ifaceItem.interface }}
              </td>
              <td>
                {{ ifaceItem.ipAddr }}
              </td>
              <td>
                {{ ifaceItem.internet }}
              </td>
              <td>
                {{ ifaceItem.result }}
              </td>
              <td>
                {{ ifaceItem.comment }}
              </td>
            </tr>
          </tbody>
        </table>

        <table class="ifaceStream">
          <colgroup>
            <col>
            <col>
            <col>
            <col>
            <col>
            <col>
            <col>
          </colgroup>
          <thead>
            <tr>
              <th>连接名称</th>
              <th>IP</th>
              <th>当前上传</th>
              <th>当前下载</th>
              <th>总上传</th>
              <th>总下载</th>
              <th>备注</th>
            </tr>
          </thead>
          <tbody v-if="item?.data?.ifaceStream && item.data.ifaceStream.length">
            <tr v-for="(ifaceStream, index) in item.data.ifaceStream" :key="index">
              <td>
                {{ ifaceStream.interface }}
              </td>
              <td>
                {{ ifaceStream.ipAddr }}
              </td>
              <td>
                {{ FormatSize(ifaceStream.upload) }}
              </td>
              <td>
                {{ FormatSize(ifaceStream.download) }}
              </td>
              <td>
                {{ FormatSize(ifaceStream.totalUp) }}
              </td>
              <td>
                {{ FormatSize(ifaceStream.totalDown) }}
              </td>
              <td>
                {{ ifaceStream.comment }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
<style scoped lang="less">
.header {
  background-color: var(--vt-c-black-mute);
  height: 56px;
  font-size: 26px;
  line-height: 56px;
  padding-left: 12px;
  font-weight: bold;
  color: #fff;
}

.interface {
  display: flex;
  justify-content: center;
  // align-items: center;
  border-radius: 2px;
  flex-wrap: wrap;

  .right {
    flex: 1;
    padding: 12px;
    color: #fff;
  }

  @media screen and (max-width: 576px) {

    .left,
    .right {
      width: 100vw;
      flex-grow: 1;
      flex-shrink: 0;
    }
  }

  @media screen and (max-width: 992px) and (min-width: 576px) {
    .left {
      width: 240px;
    }
  }

  @media screen and (min-width: 992px) {
    .left {
      width: 320px;
    }
  }

  .speed {
    .title {
      background-color: var(--vt-c-black-mute);
      height: 48px;
      font-size: 24px;
      font-weight: bold;
      background-color: rgb(24, 27, 31);
      border: 1px solid rgba(204, 204, 220, 0.12);
      padding-left: 12px;
      box-sizing: border-box;
    }

    .content {
      padding: 12px;
      border-radius: 2px;
      position: relative;
      display: flex;
      background: linear-gradient(120deg, rgb(66, 154, 67), rgb(111, 183, 87));
      align-items: center;
      justify-content: center;
      color: rgba(255, 255, 255, 0.05);
      width: 100%;
      height: 100%;

      .uploadTxt,
      .downloadTxt {
        font-weight: bold;
        text-align: center;
        display: block;
        font-size: 48px;
        color: #fff;
      }
    }

    &:last-child {
      margin-top: 12px;
    }
  }

  .interfaceGirdData {
    background: #151921;
    border-radius: 8px 8px 0 0;
  }
}
</style>
