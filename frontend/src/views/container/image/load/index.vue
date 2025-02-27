<template>
    <el-drawer v-model="loadVisiable" :destroy-on-close="true" :close-on-click-modal="false" size="30%">
        <template #header>
            <DrawerHeader :header="$t('container.importImage')" :back="handleClose" />
        </template>
        <el-form v-loading="loading" ref="formRef" :model="form" label-position="top">
            <el-row type="flex" justify="center">
                <el-col :span="22">
                    <el-form-item :label="$t('container.path')" :rules="Rules.requiredSelect" prop="path">
                        <el-input v-model="form.path">
                            <template #prepend>
                                <FileList @choose="loadLoadDir" :dir="false"></FileList>
                            </template>
                        </el-input>
                    </el-form-item>
                </el-col>
            </el-row>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button :disabeld="loading" @click="loadVisiable = false">
                    {{ $t('commons.button.cancel') }}
                </el-button>
                <el-button :disabeld="loading" type="primary" @click="onSubmit(formRef)">
                    {{ $t('container.import') }}
                </el-button>
            </span>
        </template>
    </el-drawer>
</template>

<script lang="ts" setup>
import FileList from '@/components/file-list/index.vue';
import { reactive, ref } from 'vue';
import { Rules } from '@/global/form-rules';
import i18n from '@/lang';
import { ElForm } from 'element-plus';
import { imageLoad } from '@/api/modules/container';
import { MsgSuccess } from '@/utils/message';
import DrawerHeader from '@/components/drawer-header/index.vue';

const loading = ref(false);

const loadVisiable = ref(false);
const form = reactive({
    path: '',
});

const acceptParams = () => {
    loadVisiable.value = true;
    form.path = '';
};
const handleClose = () => {
    loadVisiable.value = false;
};

const emit = defineEmits<{ (e: 'search'): void }>();

type FormInstance = InstanceType<typeof ElForm>;
const formRef = ref<FormInstance>();

const onSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    formEl.validate(async (valid) => {
        if (!valid) return;
        loading.value = true;
        await imageLoad(form)
            .then(() => {
                loading.value = false;
                loadVisiable.value = false;
                emit('search');
                MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            })
            .catch(() => {
                loading.value = false;
            });
    });
};

const loadLoadDir = async (path: string) => {
    form.path = path;
};

defineExpose({
    acceptParams,
});
</script>
