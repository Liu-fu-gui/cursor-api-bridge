<script setup>
import { listOpenAIModels } from "@/services/clientApi";
import { appState, replaceModelAdapters } from "@/state/appState";
import { computed, reactive, ref, watch } from "vue";

const props = defineProps({
  visible: { type: Boolean, default: false },
});
const emit = defineEmits(["close", "saved"]);

const form = reactive({
  baseURL: "",
  apiKey: "",
  effort: "medium",
});
const loading = ref(false);
const saving = ref(false);
const error = ref("");
const models = ref([]);
const selected = ref(new Set());
const showKey = ref(false);
const connected = computed(() => !loading.value && !error.value && models.value.length > 0);

const allSelected = computed(() => models.value.length > 0 && selected.value.size === models.value.length);

function seedFromCurrentConfig() {
  const first = appState.modelAdapters?.[0];
  if (first?.baseURL) form.baseURL = first.baseURL;
  if (first?.apiKey) form.apiKey = first.apiKey;
  if (first?.reasoningEffort) form.effort = first.reasoningEffort;
  selected.value = new Set((appState.modelAdapters || []).map((item) => item.modelID).filter(Boolean));
}

async function fetchModels() {
  error.value = "";
  const baseURL = String(form.baseURL || "").trim();
  const apiKey = String(form.apiKey || "").trim();
  if (!baseURL || !apiKey) {
    error.value = "请填写 API URL 和 Key";
    return;
  }
  loading.value = true;
  try {
    const result = await listOpenAIModels(baseURL, apiKey);
    models.value = Array.isArray(result) ? result : [];
    selected.value = new Set(models.value);
  } catch (reason) {
    error.value = String(reason?.message || reason || "模型获取失败");
    models.value = [];
  } finally {
    loading.value = false;
  }
}

function toggleModel(modelID) {
  const next = new Set(selected.value);
  if (next.has(modelID)) next.delete(modelID);
  else next.add(modelID);
  selected.value = next;
}

function toggleAll() {
  selected.value = allSelected.value ? new Set() : new Set(models.value);
}

async function save() {
  error.value = "";
  const chosen = models.value.filter((id) => selected.value.has(id));
  if (chosen.length === 0) {
    error.value = "至少选择一个模型";
    return;
  }
  saving.value = true;
  const adapters = chosen.map((modelID) => ({
    displayName: modelID,
    type: "openai",
    baseURL: String(form.baseURL || "").trim(),
    apiKey: String(form.apiKey || "").trim(),
    tooltipData: "自动获取的 OpenAI-compatible 模型",
    modelID,
    reasoningEffort: form.effort,
    openAIEndpoint: "/v1/chat/completions",
    openAIExtraParamsEnabled: false,
    openAIExtraParamsJSON: "{}",
    customHeadersEnabled: false,
    customHeadersJSON: "{}",
    contextWindowTokens: 0,
    maxCompletionTokens: 0,
  }));
  try {
    const result = await replaceModelAdapters(adapters);
    if (!result?.ok) throw new Error(result?.error || "保存失败");
    emit("saved", chosen);
    emit("close");
  } catch (reason) {
    error.value = String(reason?.message || reason || "保存失败");
  } finally {
    saving.value = false;
  }
}

watch(() => props.visible, (visible) => {
  if (!visible) return;
  showKey.value = false;
  seedFromCurrentConfig();
  void fetchModels();
});
</script>

<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-999 flex items-start justify-center bg-[#03050a]/75 px-5 pb-4 pt-12 backdrop-blur-[6px]"
      @click.self="emit('close')"
    >
      <div class="relative flex max-h-[calc(100vh-64px)] w-full max-w-[560px] flex-col overflow-hidden rounded-[20px] border border-[#252b38] bg-[#10141e] shadow-[0_30px_100px_rgba(0,0,0,.65)]">
        <div class="pointer-events-none absolute -right-20 -top-24 size-60 rounded-full bg-[#7568ff]/20 blur-3xl"></div>
        <div class="relative flex min-h-0 flex-1 flex-col p-5">
        <div class="mb-4 flex shrink-0 items-start justify-between gap-4">
          <div>
            <div class="mb-2 flex size-9 items-center justify-center rounded-[11px] bg-[#7568ff]/15 text-[#9186ff]">
              <span class="icon-[mdi--auto-fix] text-xl"></span>
            </div>
            <h3 class="text-lg font-semibold tracking-[-0.02em] text-white">连接模型服务</h3>
            <p class="mt-1 text-sm text-[#7e899f]">填写一次连接信息，自动发现并选择模型。</p>
          </div>
          <button class="flex size-8 items-center justify-center rounded-full text-xl text-[#667086] transition hover:bg-white/7 hover:text-white" @click="emit('close')">×</button>
        </div>

        <div class="min-h-0 flex-1 overflow-y-auto pr-1">
        <div class="grid gap-2.5">
          <label class="grid gap-1.5">
            <span class="text-xs font-medium text-[#9aa5bb]">API 地址</span>
            <input v-model="form.baseURL" type="text" class="h-10 rounded-[10px] border border-[#252b38] bg-[#0a0d14] px-3.5 text-sm text-white outline-none transition placeholder:text-[#424a5c] focus:border-[#7568ff]/80 focus:ring-2 focus:ring-[#7568ff]/15" />
          </label>
          <label class="grid gap-1.5">
            <span class="text-xs font-medium text-[#9aa5bb]">API Key</span>
            <div class="relative">
              <input
                v-model="form.apiKey"
                :type="showKey ? 'text' : 'password'"
                class="h-10 w-full rounded-[10px] border border-[#252b38] bg-[#0a0d14] px-3.5 pr-12 text-sm text-white outline-none transition focus:border-[#7568ff]/80 focus:ring-2 focus:ring-[#7568ff]/15"
              />
              <button type="button" class="absolute right-2 top-1/2 flex size-8 -translate-y-1/2 items-center justify-center rounded-[8px] text-[#69748a] hover:bg-white/7 hover:text-white" @click="showKey = !showKey">
                <span :class="showKey ? 'icon-[mdi--eye-off-outline]' : 'icon-[mdi--eye-outline]'"></span>
              </button>
            </div>
          </label>
          <button
            class="mt-0.5 flex h-10 items-center justify-center gap-2 rounded-[10px] border border-[#7568ff]/35 bg-[#7568ff]/12 text-sm font-medium text-[#aaa3ff] transition hover:border-[#7568ff]/60 hover:bg-[#7568ff]/20 disabled:opacity-55"
            :disabled="loading"
            @click="fetchModels"
          >
            <span :class="loading ? 'icon-[mdi--loading] animate-spin' : 'icon-[mdi--radar]'"></span>
            {{ loading ? "正在连接并发现模型…" : "测试连接并获取模型" }}
          </button>
        </div>

        <div v-if="error" class="mt-3 rounded-[10px] border border-[#ff6575]/25 bg-[#ff4455]/10 px-3 py-2.5 text-sm text-[#ffabb3]">{{ error }}</div>

        <div class="mt-3 overflow-hidden rounded-[12px] border border-[#252b38] bg-[#0b0e15]">
          <div class="flex items-center justify-between border-b border-white/8 px-3.5 py-3">
            <div class="flex items-center gap-2">
              <span class="size-2 rounded-full" :class="connected ? 'bg-[#37e6a5] shadow-[0_0_10px_#37e6a5]' : 'bg-[#4a5263]'"></span>
              <span class="text-sm text-[#cbd2e1]">{{ connected ? `连接成功 · ${models.length} 个模型` : "等待连接" }}</span>
            </div>
            <button v-if="models.length" class="text-xs font-medium text-[#9186ff] hover:text-[#aaa3ff]" @click="toggleAll">{{ allSelected ? "取消全选" : "全部选择" }}</button>
          </div>
          <div class="grid max-h-[130px] grid-cols-2 gap-1.5 overflow-y-auto p-2.5">
            <label
              v-for="(modelID, index) in models"
              :key="modelID"
              class="flex cursor-pointer items-center gap-2.5 rounded-[9px] border px-3 py-2 transition"
              :class="[
                selected.has(modelID) ? 'border-[#7568ff]/35 bg-[#7568ff]/12' : 'border-[#1c212c] bg-[#11151e] hover:bg-[#171c27]',
                models.length % 2 === 1 && index === models.length - 1 ? 'col-span-2' : ''
              ]"
            >
              <input type="checkbox" :checked="selected.has(modelID)" class="size-4 accent-[#7568ff]" @change="toggleModel(modelID)" />
              <span class="min-w-0 truncate text-sm" :class="selected.has(modelID) ? 'text-white' : 'text-[#9aa5bb]'">{{ modelID }}</span>
            </label>
            <div v-if="!loading && models.length === 0" class="col-span-2 px-3 py-7 text-center text-sm text-[#515a6d]">输入连接信息后自动发现模型</div>
          </div>
        </div>

        </div>

        <div class="sticky bottom-0 mt-3 flex shrink-0 items-end justify-between gap-4 border-t border-[#202633] bg-[#10141e] pt-3">
          <label class="grid min-w-[170px] gap-1.5">
            <span class="text-xs font-medium text-[#9aa5bb]">默认推理强度</span>
            <select v-model="form.effort" class="h-10 rounded-[10px] border border-[#252b38] bg-[#0a0d14] px-3 text-sm text-white outline-none focus:border-[#7568ff]/70">
              <option value="low">轻度</option>
              <option value="medium">中等</option>
              <option value="high">高</option>
              <option value="xhigh">极高</option>
              <option value="max">最大</option>
            </select>
          </label>
          <div class="flex gap-2">
            <button class="h-10 rounded-[10px] px-4 text-sm text-[#8490a7] transition hover:bg-white/6 hover:text-white" @click="emit('close')">取消</button>
            <button
              class="h-10 min-w-[120px] rounded-[10px] bg-gradient-to-r from-[#7568ff] to-[#5c8dff] px-4 text-sm font-medium text-white shadow-[0_8px_24px_rgba(105,106,255,.3)] transition hover:-translate-y-0.5 disabled:translate-y-0 disabled:opacity-45"
              :disabled="saving || selected.size === 0"
              @click="save"
            >
              {{ saving ? "正在应用…" : "保存并应用" }}
            </button>
          </div>
        </div>
        </div>
        </div>
      </div>
  </Teleport>
</template>
