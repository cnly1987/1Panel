<template>
    <el-drawer v-model="drawerVisiable" :destroy-on-close="true" :close-on-click-modal="false" size="50%">
        <template #header>
            <DrawerHeader :header="$t('container.compose')" :back="handleClose" />
        </template>
        <div v-loading="loading">
            <el-form ref="formRef" label-position="top" :model="form" :rules="rules" label-width="80px">
                <el-row type="flex" justify="center">
                    <el-col :span="22">
                        <el-form-item :label="$t('container.from')">
                            <el-radio-group v-model="form.from">
                                <el-radio label="edit">{{ $t('commons.button.edit') }}</el-radio>
                                <el-radio label="path">{{ $t('container.pathSelect') }}</el-radio>
                                <el-radio label="template">{{ $t('container.composeTemplate') }}</el-radio>
                            </el-radio-group>
                        </el-form-item>
                        <el-form-item v-if="form.from === 'path'" prop="path">
                            <el-input
                                :placeholder="$t('commons.example') + '/tmp/docker-compose.yml'"
                                v-model="form.path"
                            >
                                <template #prepend>
                                    <FileList @choose="loadDir" :dir="false"></FileList>
                                </template>
                            </el-input>
                        </el-form-item>
                        <el-form-item v-if="form.from === 'edit' || form.from === 'template'" prop="name">
                            <el-input @input="changePath" v-model.trim="form.name">
                                <template #prepend>{{ $t('file.dir') }}</template>
                            </el-input>
                            <span class="input-help">{{ $t('container.composePathHelper', [composeFile]) }}</span>
                        </el-form-item>
                        <el-form-item v-if="form.from === 'template'" prop="template">
                            <el-select v-model="form.template">
                                <el-option
                                    v-for="item in templateOptions"
                                    :key="item.id"
                                    :value="item.id"
                                    :label="item.name"
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item v-if="form.from === 'edit'">
                            <codemirror
                                :autofocus="true"
                                placeholder="#Define or paste the content of your docker-compose file here"
                                :indent-with-tab="true"
                                :tabSize="4"
                                style="width: 100%; height: calc(100vh - 340px)"
                                :lineWrapping="true"
                                :matchBrackets="true"
                                theme="cobalt"
                                :styleActiveLine="true"
                                :extensions="extensions"
                                v-model="form.file"
                            />
                        </el-form-item>
                    </el-col>
                </el-row>
            </el-form>
        </div>
        <template #footer>
            <span class="dialog-footer">
                <el-button :disabled="loading" @click="drawerVisiable = false">
                    {{ $t('commons.button.cancel') }}
                </el-button>
                <el-button type="primary" :disabled="loading" @click="onSubmit(formRef)">
                    {{ $t('commons.button.confirm') }}
                </el-button>
            </span>
        </template>
    </el-drawer>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import FileList from '@/components/file-list/index.vue';
import { Codemirror } from 'vue-codemirror';
import { javascript } from '@codemirror/lang-javascript';
import { oneDark } from '@codemirror/theme-one-dark';
import { Rules } from '@/global/form-rules';
import i18n from '@/lang';
import { ElForm } from 'element-plus';
import DrawerHeader from '@/components/drawer-header/index.vue';
import { listComposeTemplate, upCompose } from '@/api/modules/container';
import { MsgSuccess } from '@/utils/message';
import { loadBaseDir } from '@/api/modules/setting';

const extensions = [javascript(), oneDark];
const drawerVisiable = ref(false);
const templateOptions = ref();

const loading = ref(false);
const baseDir = ref();
const composeFile = ref();

const varifyPath = (rule: any, value: any, callback: any) => {
    if (value.indexOf('docker-compose.yml') === -1) {
        callback(new Error(i18n.global.t('commons.rule.selectHelper', ['docker-compose.yml'])));
    }
    callback();
};
const form = reactive({
    name: '',
    from: 'edit',
    path: '',
    file: '',
    template: null as number,
});
const rules = reactive({
    name: [Rules.requiredInput, Rules.imageName],
    path: [Rules.requiredSelect, { validator: varifyPath, trigger: 'change', required: true }],
});

const loadTemplates = async () => {
    const res = await listComposeTemplate();
    templateOptions.value = res.data;
    if (templateOptions.value && templateOptions.value.length !== 0) {
        form.template = templateOptions.value[0].id;
    }
};

const acceptParams = (): void => {
    drawerVisiable.value = true;
    form.name = '';
    form.from = 'edit';
    form.path = '';
    form.file = '';
    loadTemplates();
    loadPath();
};
const emit = defineEmits<{ (e: 'search'): void }>();

const handleClose = () => {
    drawerVisiable.value = false;
};

const loadPath = async () => {
    const pathRes = await loadBaseDir();
    baseDir.value = pathRes.data;
    changePath();
};

const changePath = async () => {
    composeFile.value = baseDir.value + '/docker/compose/' + form.name + '/docker-compose.yml';
};

type FormInstance = InstanceType<typeof ElForm>;
const formRef = ref<FormInstance>();

const onSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    formEl.validate(async (valid) => {
        if (!valid) return;
        loading.value = true;
        upCompose(form)
            .then(() => {
                MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
                loading.value = false;
                emit('search');
                drawerVisiable.value = false;
            })
            .finally(() => {
                loading.value = false;
            });
    });
};

const loadDir = async (path: string) => {
    form.path = path;
};

defineExpose({
    acceptParams,
});
</script>
