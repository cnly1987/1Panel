<template>
    <LayoutContent v-loading="loading || syncLoading" :title="activeName" :divider="true">
        <template #toolbar>
            <el-row :gutter="5">
                <el-col :span="20">
                    <div>
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
                                {{ item.name }}
                            </el-button>
                        </div>
                    </div>
                </el-col>
                <el-col :span="4">
                    <div class="search-button">
                        <el-input
                            class="table-button"
                            v-model="searchReq.name"
                            clearable
                            @clear="search()"
                            suffix-icon="Search"
                            @keyup.enter="search()"
                            @blur="search()"
                            :placeholder="$t('commons.button.search')"
                        ></el-input>
                    </div>
                </el-col>
            </el-row>
        </template>
        <template #rightButton>
            <el-button @click="sync" type="primary" link v-if="mode === 'installed' && data != null">
                {{ $t('app.sync') }}
            </el-button>
        </template>

        <template #main>
            <div class="update-prompt" v-if="data == null">
                <span>{{ mode === 'upgrade' ? $t('app.updatePrompt') : $t('app.installPrompt') }}</span>
                <div>
                    <img src="@/assets/images/no_update_app.svg" />
                </div>
            </div>
            <el-row :gutter="5">
                <el-col
                    v-for="(installed, index) in data"
                    :key="index"
                    :xs="24"
                    :sm="24"
                    :md="24"
                    :lg="12"
                    :xl="12"
                    class="install-card-col-12"
                >
                    <div class="install-card">
                        <el-row :gutter="24">
                            <el-col :xs="3" :sm="3" :md="3" :lg="4" :xl="4">
                                <div class="icon">
                                    <el-avatar
                                        shape="square"
                                        :size="66"
                                        :src="'data:image/png;base64,' + installed.app.icon"
                                    />
                                </div>
                            </el-col>
                            <el-col :xs="21" :sm="21" :md="21" :lg="20" :xl="20">
                                <div class="a-detail">
                                    <div class="d-name">
                                        <span class="name">{{ installed.name }}</span>
                                        <span class="status">
                                            <el-popover
                                                v-if="installed.status === 'Error'"
                                                placement="bottom"
                                                :width="400"
                                                trigger="hover"
                                                :content="installed.message"
                                            >
                                                <template #reference>
                                                    <Status :key="installed.status" :status="installed.status"></Status>
                                                </template>
                                            </el-popover>
                                            <span v-else>
                                                <Status :key="installed.status" :status="installed.status"></Status>
                                            </span>
                                        </span>

                                        <el-button
                                            class="h-button"
                                            type="primary"
                                            plain
                                            round
                                            size="small"
                                            :disabled="installed.status !== 'Running'"
                                            @click="openUploads(installed.app.key, installed.name)"
                                            v-if="mode === 'installed'"
                                        >
                                            {{ $t('database.loadBackup') }}
                                        </el-button>
                                        <el-button
                                            class="h-button"
                                            type="primary"
                                            plain
                                            round
                                            size="small"
                                            :disabled="installed.status !== 'Running'"
                                            @click="openBackups(installed.app.key, installed.name)"
                                            v-if="mode === 'installed'"
                                        >
                                            {{ $t('app.backup') }}
                                        </el-button>
                                        <el-button
                                            class="h-button"
                                            type="primary"
                                            plain
                                            round
                                            size="small"
                                            @click="openOperate(installed, 'upgrade')"
                                            v-if="mode === 'upgrade'"
                                        >
                                            {{ $t('app.upgrade') }}
                                        </el-button>
                                    </div>
                                    <div class="d-description">
                                        <el-tag>{{ $t('app.version') }}：{{ installed.version }}</el-tag>
                                        <el-tag v-if="installed.httpPort > 0">
                                            {{ $t('app.busPort') }}：{{ installed.httpPort }}
                                        </el-tag>
                                        <div class="description">
                                            <span>{{ $t('app.areadyRun') }}： {{ getAge(installed.createdAt) }}</span>
                                        </div>
                                    </div>
                                    <div class="app-divider" />
                                    <div
                                        class="d-button"
                                        v-if="mode === 'installed' && installed.status != 'Installing'"
                                    >
                                        <el-button
                                            v-for="(button, key) in buttons"
                                            :key="key"
                                            :type="button.disabled && button.disabled(installed) ? 'info' : 'primary'"
                                            plain
                                            round
                                            size="small"
                                            @click="button.click(installed)"
                                            :disabled="button.disabled && button.disabled(installed)"
                                        >
                                            {{ button.label }}
                                        </el-button>
                                    </div>
                                </div>
                            </el-col>
                        </el-row>
                    </div>
                </el-col>
            </el-row>
        </template>
    </LayoutContent>
    <Backups ref="backupRef" @close="search" />
    <Uploads ref="uploadRef" />
    <AppResources ref="checkRef" />
    <AppDelete ref="deleteRef" @close="search" />
    <AppParams ref="appParamRef" />
    <AppUpgrade ref="upgradeRef" @close="search" />
</template>

<script lang="ts" setup>
import {
    SearchAppInstalled,
    InstalledOp,
    SyncInstalledApp,
    AppInstalledDeleteCheck,
    GetAppTags,
} from '@/api/modules/app';
import LayoutContent from '@/layout/layout-content.vue';
import { onMounted, onUnmounted, reactive, ref } from 'vue';
import i18n from '@/lang';
import { ElMessageBox } from 'element-plus';
import Backups from '@/components/backup/index.vue';
import Uploads from '@/components/upload/index.vue';
import AppResources from './check/index.vue';
import AppDelete from './delete/index.vue';
import AppParams from './detail/index.vue';
import AppUpgrade from './upgrade/index.vue';
import { App } from '@/api/interface/app';
import Status from '@/components/status/index.vue';
import { getAge } from '@/utils/util';
import { useRouter } from 'vue-router';
import { MsgSuccess } from '@/utils/message';

let data = ref<any>();
let loading = ref(false);
let syncLoading = ref(false);
let timer: NodeJS.Timer | null = null;
const paginationConfig = reactive({
    currentPage: 1,
    pageSize: 20,
    total: 0,
});
let open = ref(false);
let operateReq = reactive({
    installId: 0,
    operate: '',
    detailId: 0,
});
const backupRef = ref();
const uploadRef = ref();
const checkRef = ref();
const deleteRef = ref();
const appParamRef = ref();
const upgradeRef = ref();
let tags = ref<App.Tag[]>([]);
let activeTag = ref('all');
let searchReq = reactive({
    page: 1,
    pageSize: 15,
    name: '',
    tags: [],
    update: false,
});
const router = useRouter();
let activeName = ref(i18n.global.t('app.installed'));
let mode = ref('installed');

const sync = () => {
    syncLoading.value = true;
    SyncInstalledApp()
        .then(() => {
            MsgSuccess(i18n.global.t('app.syncSuccess'));
            search();
        })
        .finally(() => {
            syncLoading.value = false;
        });
};

const changeTag = (key: string) => {
    searchReq.tags = [];
    activeTag.value = key;
    if (key !== 'all') {
        searchReq.tags = [key];
    }
    search();
};

const search = () => {
    loading.value = true;
    searchReq.page = paginationConfig.currentPage;
    searchReq.pageSize = paginationConfig.pageSize;
    SearchAppInstalled(searchReq)
        .then((res) => {
            data.value = res.data.items;
            paginationConfig.total = res.data.total;
        })
        .finally(() => {
            loading.value = false;
        });
    GetAppTags().then((res) => {
        tags.value = res.data;
    });
};

const openOperate = (row: any, op: string) => {
    operateReq.installId = row.id;
    operateReq.operate = op;
    if (op == 'upgrade') {
        upgradeRef.value.acceptParams(row.id, row.name);
    } else if (op == 'delete') {
        AppInstalledDeleteCheck(row.id).then(async (res) => {
            const items = res.data;
            if (res.data && res.data.length > 0) {
                checkRef.value.acceptParams({ items: items });
            } else {
                deleteRef.value.acceptParams(row);
            }
        });
    } else {
        onOperate(op);
    }
};

const operate = async () => {
    open.value = false;
    loading.value = true;
    await InstalledOp(operateReq)
        .then(() => {
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            search();
        })
        .catch(() => {
            search();
        })
        .finally(() => {
            loading.value = false;
        });
};

const onOperate = async (operation: string) => {
    ElMessageBox.confirm(
        i18n.global.t('app.operatorHelper', [i18n.global.t('app.' + operation)]),
        i18n.global.t('app.' + operation),
        {
            confirmButtonText: i18n.global.t('commons.button.confirm'),
            cancelButtonText: i18n.global.t('commons.button.cancel'),
            type: 'info',
        },
    ).then(() => {
        operate();
    });
};

const buttons = [
    {
        label: i18n.global.t('app.sync'),
        click: (row: any) => {
            openOperate(row, 'sync');
        },
    },
    {
        label: i18n.global.t('app.rebuild'),
        click: (row: any) => {
            openOperate(row, 'rebuild');
        },
    },
    {
        label: i18n.global.t('app.restart'),
        click: (row: any) => {
            openOperate(row, 'restart');
        },
    },
    {
        label: i18n.global.t('app.start'),
        click: (row: any) => {
            openOperate(row, 'start');
        },
        disabled: (row: any) => {
            return row.status === 'Running';
        },
    },
    {
        label: i18n.global.t('app.stop'),
        click: (row: any) => {
            openOperate(row, 'stop');
        },
        disabled: (row: any) => {
            return row.status !== 'Running';
        },
    },
    {
        label: i18n.global.t('app.delete'),
        click: (row: any) => {
            openOperate(row, 'delete');
        },
    },
    {
        label: i18n.global.t('app.params'),
        click: (row: any) => {
            openParam(row.id);
        },
    },
];

const openBackups = (key: string, name: string) => {
    let params = {
        type: 'app',
        name: key,
        detailName: name,
    };
    backupRef.value.acceptParams(params);
};

const openUploads = (key: string, name: string) => {
    let params = {
        type: 'app',
        name: key,
        detailName: name,
    };
    uploadRef.value.acceptParams(params);
};

const openParam = (installId: number) => {
    appParamRef.value.acceptParams({ id: installId });
};

onMounted(() => {
    const path = router.currentRoute.value.path;
    if (path == '/apps/upgrade') {
        activeName.value = i18n.global.t('app.canUpgrade');
        mode.value = 'upgrade';
        searchReq.update = true;
    }
    search();
    timer = setInterval(() => {
        search();
    }, 10000 * 6);
});

onUnmounted(() => {
    clearInterval(Number(timer));
    timer = null;
});
</script>

<style scoped lang="scss">
@import '../index.scss';
@media only screen and (max-width: 1300px) {
    .install-card-col-12 {
        max-width: 100%;
        flex: 0 0 100%;
    }
}
</style>
