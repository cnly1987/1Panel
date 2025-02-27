<template>
    <LayoutContent v-loading="loading" v-if="!showDetail" :title="$t('app.app')" :divider="true">
        <template #toolbar>
            <el-row :gutter="5">
                <el-col :span="20">
                    <el-button
                        class="tag-button"
                        :class="activeTag === 'all' ? '' : 'no-active'"
                        @click="changeTag('all')"
                        :type="activeTag === 'all' ? 'primary' : ''"
                        :plain="activeTag !== 'all'"
                    >
                        {{ $t('app.all') }}
                    </el-button>
                    <div v-for="item in tags" :key="item.key" style="display: inline">
                        <el-button
                            class="tag-button"
                            :class="activeTag === item.key ? '' : 'no-active'"
                            @click="changeTag(item.key)"
                            :type="activeTag === item.key ? 'primary' : ''"
                            :plain="activeTag !== item.key"
                        >
                            {{ language == 'zh' ? item.name : item.key }}
                        </el-button>
                    </div>
                </el-col>
                <el-col :span="4">
                    <div class="search-button">
                        <el-input
                            v-model="req.name"
                            clearable
                            @clear="searchByName('')"
                            suffix-icon="Search"
                            @keyup.enter="searchByName(req.name)"
                            @blur="searchByName(req.name)"
                            :placeholder="$t('commons.button.search')"
                        ></el-input>
                    </div>
                </el-col>
            </el-row>
        </template>
        <template #rightButton>
            <el-badge is-dot class="item" :hidden="!canUpdate">
                <el-button @click="sync" type="primary" link :plain="true">{{ $t('app.syncAppList') }}</el-button>
            </el-badge>
        </template>
        <template #main>
            <el-row :gutter="5">
                <el-col v-for="(app, index) in apps" :key="index" :xs="12" :sm="12" :md="8" :lg="8" :xl="8">
                    <div class="app-card">
                        <el-row :gutter="24">
                            <el-col :xs="5" :sm="5" :md="6" :lg="6" :xl="5">
                                <div class="app-icon">
                                    <el-avatar shape="square" :size="60" :src="'data:image/png;base64,' + app.icon" />
                                </div>
                            </el-col>
                            <el-col :xs="19" :sm="19" :md="18" :lg="18" :xl="19">
                                <div class="app-content">
                                    <div class="app-header">
                                        <span class="app-title">{{ app.name }}</span>
                                        <el-button
                                            class="app-button"
                                            type="primary"
                                            plain
                                            round
                                            size="small"
                                            @click="getAppDetail(app.key)"
                                        >
                                            {{ $t('app.install') }}
                                        </el-button>
                                    </div>
                                    <div class="app-desc">
                                        <span class="desc">
                                            {{ language == 'zh' ? app.shortDescZh : app.shortDescEn }}
                                        </span>
                                    </div>
                                    <div class="app-tag">
                                        <el-tag v-for="(tag, ind) in app.tags" :key="ind" :colr="getColor(ind)">
                                            {{ language == 'zh' ? tag.name : tag.key }}
                                        </el-tag>
                                    </div>
                                    <div class="divider"></div>
                                </div>
                            </el-col>
                        </el-row>
                    </div>
                </el-col>
            </el-row>
        </template>
    </LayoutContent>
    <Detail v-if="showDetail" :id="appId"></Detail>
</template>

<script lang="ts" setup>
import LayoutContent from '@/layout/layout-content.vue';
import { App } from '@/api/interface/app';
import { onMounted, reactive, ref } from 'vue';
import { GetAppListUpdate, GetAppTags, SearchApp, SyncApp } from '@/api/modules/app';
import i18n from '@/lang';
import Detail from '../detail/index.vue';
import router from '@/routers';
import { MsgSuccess } from '@/utils/message';
import { useI18n } from 'vue-i18n';

const language = useI18n().locale.value;

let req = reactive({
    name: '',
    tags: [],
    page: 1,
    pageSize: 50,
});

let apps = ref<App.App[]>([]);
let tags = ref<App.Tag[]>([]);
const colorArr = ['#6495ED', '#54FF9F', '#BEBEBE', '#FFF68F', '#FFFF00', '#8B0000'];
let loading = ref(false);
let activeTag = ref('all');
let showDetail = ref(false);
let appId = ref(0);
let canUpdate = ref(false);

const getColor = (index: number) => {
    return colorArr[index];
};

const search = async (req: App.AppReq) => {
    loading.value = true;
    await SearchApp(req)
        .then((res) => {
            apps.value = res.data.items;
        })
        .finally(() => {
            loading.value = false;
        });
    GetAppTags().then((res) => {
        tags.value = res.data;
    });
};

const getAppDetail = (key: string) => {
    router.push({ name: 'AppDetail', params: { appKey: key } });
};

const sync = () => {
    loading.value = true;
    SyncApp()
        .then(() => {
            MsgSuccess(i18n.global.t('app.syncSuccess'));
            canUpdate.value = false;
            search(req);
        })
        .finally(() => {
            loading.value = false;
        });
};

const getAppListUpdate = async () => {
    const res = await GetAppListUpdate();
    canUpdate.value = res.data.canUpdate;
};

const changeTag = (key: string) => {
    req.tags = [];
    activeTag.value = key;
    if (key !== 'all') {
        req.tags = [key];
    }
    search(req);
};

const searchByName = (name: string) => {
    req.name = name;
    search(req);
};

onMounted(() => {
    getAppListUpdate();
    search(req);
});
</script>

<style lang="scss">
.header {
    padding-bottom: 10px;
}

.app-card {
    margin-top: 10px;
    cursor: pointer;
    padding: 5px;

    .app-icon {
        margin-top: 10px;
        margin-left: 10px;
    }

    .app-content {
        margin-top: 10px;
        height: 100%;
        width: 100%;

        .app-header {
            height: 20%;
            .app-title {
                font-weight: 500;
                font-size: 16px;
                color: var(--el-text-color-regular);
            }
            .app-button {
                float: right;
                margin-right: 20px;
            }
        }

        .app-desc {
            margin-top: 5px;
            overflow: hidden;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;

            text-overflow: ellipsis;
            height: 45px;

            .desc {
                font-size: 14px;
                color: var(--el-text-color-regular);
            }
        }

        .app-tag {
            margin-top: 5px;
        }
    }
    &:hover {
        background-color: rgba(0, 94, 235, 0.03);
    }
}

.tag-button {
    margin-right: 10px;
    &.no-active {
        background: none;
        border: none;
    }
}
</style>
