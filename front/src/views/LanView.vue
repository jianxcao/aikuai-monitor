<script setup lang="ts">
import { normalizeKey } from '@/utils';
import { computed, onMounted, ref } from 'vue';
import { FormatSize, FormatNum } from '@/utils'
import { useUrlSearchParams } from '@vueuse/core'
import IconUp from '@/components/icons/IconUp.vue';
import IconDown from '@/components/icons/IconDown.vue';

const params = useUrlSearchParams('history')
const lanv6 = ref<Record<string, any>>({})
const lanv4 = ref<Record<string, any>>({})
const sortBy = ref<Record<string, {
    columnName?: string
    sort?: number
}>>({})
const filterBy = ref<Record<string, string[]>>({})
let timeId: number | undefined
const headers = [
    { key: "mac", title: "Mac", },
    { key: "ipAddr", title: "ip", },
    { key: "upload", title: "当前上传", },
    { key: "download", title: "当前下载", },
    { key: "totalUp", title: "总上传", },
    { key: "totalDown", title: "总下载", },
    { key: "connectNum", title: "连接数", },
]

try {
    const filter = JSON.parse(params.filter as string || '{}')
    filterBy.value = filter
} catch (error) {
    console.error(error)
}

const data = computed(() => {
    const ipToNumber = (ip: string) =>
        ip
            .split(".")
            .map((octet: string) => parseInt(octet, 10))
            .reduce((acc, val) => (acc << 8) + val);
    return Object.keys(lanv4.value).reduce((prev: Record<string, any>, key) => {
        const ipv4 = lanv4.value[key]?.data || { total: 0, data: [] }
        const mapIpv6Data = (lanv6.value[key]?.data?.data || []).reduce((prev: any, cur: any) => {
            if (cur && cur.mac) {
                prev[cur.mac] = cur;
            }
            return prev;
        }, {})
        let data = ipv4.data.map((item: any) => {
            const mac = item.mac
            const v6Lan = mapIpv6Data[mac]
            delete mapIpv6Data[mac]
            const result = {
                v4: item,
                v6: v6Lan,
            }
            return result;
        })
        if (Object.keys(mapIpv6Data).length > 0) {
            for (let key in mapIpv6Data) {
                data.push({
                    v4: undefined,
                    v6: mapIpv6Data[key]
                })
            }
        }

        data = data.filter((item: any) => {
            if (filterBy.value[key]?.length) {
                return filterBy.value[key].some(str => str ? item.v4?.comment.includes(str) : true) ||
                    filterBy.value[key].some(str => str ? item.v6?.comment.includes(str) : true)
            }
            return true
        })

        if (sortBy.value[key] && Number(sortBy.value[key].sort) >= 0) {
            const columnName = sortBy.value[key].columnName || ""
            const sortType = sortBy.value[key].sort
            data = data.sort((a: any, b: any) => {
                if (columnName === 'mac') {
                    let aMac = (a.v4?.mac || a.v6?.mac)
                    let bMac = (b.v4?.mac || b.v6?.mac)
                    if (aMac < bMac) {
                        return sortType === 1 ? 1 : -1;
                    } else {
                        return sortType === 1 ? -1 : 0;
                    }
                } else if (columnName === 'ipAddr') {
                    // 只排ipv4，其他的都放到最后
                    if (a.v4 || b.v4) {
                        if (a.v4 && b.v4) {
                            const aip = ipToNumber(a.v4.ipAddr)
                            const bip = ipToNumber(b.v4.ipAddr)
                            if (aip < bip) {
                                return sortType === 1 ? 1 : -1;
                            } else {
                                return sortType === 1 ? -1 : 0;
                            }
                        } else if (a.v4) {
                            return -1
                        } else if (b.v4) {
                            return 1
                        }
                    }
                    return 1;
                } else if (["upload", "download", "totalUp", "totalDown", "connectNum"].includes(columnName)) {
                    let aField = (a.v4 ? Number(a.v4[columnName]) : 0) + (a.v6 ? Number(a.v6[columnName]) : 0);
                    let bField = b.v4 ? Number(b.v4[columnName]) : 0 + (b.v6 ? Number(b.v6[columnName]) : 0);
                    if (aField < bField) {
                        return sortType === 1 ? 1 : -1;
                    } else {
                        return sortType === 1 ? -1 : 0;
                    }
                }
                return 0;
            })
        }

        prev[key] = {
            total: ipv4.total,
            data
        }
        return prev;
    }, {})
})

const fetchLan = async (isV6 = false) => {
    try {
        const res = await fetch(`/api/lanv${isV6 ? '6' : '4'}`)
        let d = await res.json()
        d = normalizeKey.underlineToCamelCase(d) || {};
        return d.data || {};
    } catch (error) {
        console.error(error)
        return {}
    }
}

const fetchAll = async () => {
    const reuslt = await Promise.all([fetchLan(), fetchLan(true)])
    lanv4.value = reuslt[0];
    lanv6.value = reuslt[1];
}

const setFitler = (key: string, evt: KeyboardEvent) => {
    const inputElement = evt.target as HTMLInputElement;
    filterBy.value[key] = inputElement.value.split(',')
    params.filter = JSON.stringify(filterBy.value)
}

const sort = (key: string, column: string) => {
    sortBy.value[key] = sortBy.value[key] || {}
    sortBy.value[key].columnName = column
    const s = sortBy.value[key].sort
    if (s === undefined || s === 1) {
        sortBy.value[key].sort = 0;
    } else {
        sortBy.value[key].sort = 1;
    }
    console.debug(sortBy)
}

onMounted(async () => {
    await fetchAll()
    timeId = setInterval(() => {
        fetchAll()
    }, 3000)
})

onMounted(() => {
    clearInterval(timeId)
})
</script>


<template>
    <div class="lan-wrapper" v-for="(item, key) in data" :key="key">
        <div class="header">
            <span>{{ key }}</span> <span><input placeholder="通过备注过滤" @keyup.enter="setFitler(key, $event)"
                    :value="filterBy[key]?.join(',')"></span>
        </div>
        <table class="lan-table">
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
                    <th v-for="(headerItem, idx) in headers" :key="idx" @click="sort(key, headerItem.key)">
                        <div class="header-cell">
                            <span> {{ headerItem.title }}</span>
                            <span v-if="sortBy[key] && sortBy[key].columnName === headerItem.key">
                                <IconUp v-if="sortBy[key].sort === 0" />
                                <IconDown v-if="sortBy[key].sort === 1" />
                            </span>
                        </div>
                    </th>
                    <th>备注</th>
                </tr>
            </thead>
            <tbody v-if="item?.data && item.data.length">
                <tr v-for="(lanItem, index) in item.data" :key="index">
                    <td>
                        {{ lanItem.v4?.mac || lanItem.v6?.mac }}
                    </td>
                    <td>
                        <div>{{ lanItem.v4?.ipAddr }}</div>
                        <div>{{ lanItem.v6?.ipAddr }}</div>
                    </td>
                    <td>
                        <div>{{ FormatSize(lanItem.v4?.upload) }}</div>
                        <div> {{ FormatSize(lanItem.v6?.upload) }}</div>
                    </td>
                    <td>
                        <div> {{ FormatSize(lanItem.v4?.download) }}</div>
                        <div>{{ FormatSize(lanItem.v6?.download) }}</div>
                    </td>
                    <td>
                        <div>{{ FormatSize(lanItem.v4?.totalUp) }}</div>
                        <div>{{ FormatSize(lanItem.v6?.totalUp) }}</div>
                    </td>
                    <td>
                        <div> {{ FormatSize(lanItem.v4?.totalDown) }}</div>
                        <div>{{ FormatSize(lanItem.v6?.totalDown) }}</div>
                    </td>
                    <td>
                        <div>{{ FormatNum(lanItem.v4?.connectNum) }}</div>
                        <div>{{ FormatNum(lanItem.v6?.connectNum) }}</div>
                    </td>
                    <td>
                        {{ lanItem.v4?.comment || lanItem.v6?.comment }}
                    </td>
                </tr>
            </tbody>
        </table>
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
    display: flex;
    align-items: center;

    >span {
        display: block;
        padding-right: 12px;
        display: inline-flex;
        align-items: center;

        input {
            height: 32px;
            line-height: 32px;
            background-color: #313136;
            border: 1px solid #10b981;
            border-radius: 4px;
            font-size: 14px;
            color: rgba(255, 255, 245, .86);
            padding: 0 8px;
            outline: none;

            &:active {
                border: 1px solid #10b981;
            }
        }
    }
}

.lan-table {
    .header-cell {
        display: flex;
        align-items: center;

        &>span {
            display: inline-block;
            margin-right: 4px;
            height: 24px;
            line-height: 24px;
        }
    }
}
</style>