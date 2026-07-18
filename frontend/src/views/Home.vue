<script setup>
import QuickModelSetupModal from "@/components/QuickModelSetupModal.vue";
import { useMessage } from "@/composables/useMessage";
import { showModal } from "@/composables/useModal";
import {
  appState,
  appViewState,
  openConfigWindow,
  syncHomeMetrics,
  toUserError,
  toggleService,
} from "@/state/appState";
import { computed, onMounted, ref } from "vue";

const message = useMessage();
const quickModelSetupVisible = ref(false);
const configuredModels = computed(() => appState.modelAdapters?.length || 0);
const activeURL = computed(() => appState.modelAdapters?.[0]?.baseURL || "尚未配置");
const tokenTotal = computed(() => appState.homeMetrics?.requestTokensTotal || appState.homeMetrics?.tokensTotal || 0);

async function showActionError(title, error) {
  await showModal({
    title,
    content: String(error || "服务错误").trim() || "服务错误",
  });
}

async function handleToggleService() {
  const result = await toggleService();
  if (!result.ok) await showActionError("服务操作失败", result.error);
}

async function handleOpenConfig() {
  try {
    await openConfigWindow();
  } catch (error) {
    await showActionError("打开失败", toUserError(error));
  }
}

onMounted(() => {
  void syncHomeMetrics();
});
</script>

<template>
  <main class="min-h-full overflow-y-auto bg-[#080b12] px-5 pb-6 text-white">
    <section class="relative overflow-hidden rounded-[22px] border border-white/10 bg-gradient-to-br from-[#17213d] via-[#10182c] to-[#0b1020] p-6 shadow-[0_24px_80px_rgba(0,0,0,.35)]">
      <div class="pointer-events-none absolute -right-20 -top-24 size-72 rounded-full bg-[#6d5dfc]/25 blur-3xl"></div>
      <div class="pointer-events-none absolute -bottom-24 left-16 size-64 rounded-full bg-[#19c7a3]/15 blur-3xl"></div>

      <div class="relative flex items-start justify-between gap-5">
        <div>
          <div class="mb-3 inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3 py-1.5 text-xs text-[#b9c4df]">
            <span class="size-2 rounded-full" :class="appState.serviceRunning ? 'bg-[#37e6a5] shadow-[0_0_12px_#37e6a5]' : 'bg-[#687086]'"></span>
            {{ appViewState.serviceStatusText }}
          </div>
          <h1 class="text-[28px] font-semibold tracking-[-0.03em]">Cursor API 转换器</h1>
          <p class="mt-2 max-w-[430px] text-sm leading-6 text-[#9ca9c5]">
            让 Cursor Agent 直接使用你的模型服务。配置保存在本机，不依赖 Cursor 账号额度。
          </p>
        </div>

        <button
          class="group flex min-w-[132px] items-center justify-center gap-2 rounded-[14px] px-5 py-3 text-sm font-medium transition-all disabled:opacity-50"
          :class="appState.serviceRunning
            ? 'border border-white/10 bg-white/7 text-[#d8def0] hover:bg-white/12'
            : 'bg-gradient-to-r from-[#7568ff] to-[#5a8cff] text-white shadow-[0_10px_30px_rgba(91,105,255,.35)] hover:-translate-y-0.5'"
          :disabled="appState.serviceBusy"
          @click="handleToggleService"
        >
          <span :class="appState.serviceRunning ? 'icon-[mdi--stop-circle-outline]' : 'icon-[mdi--play-circle-outline]'" class="text-lg"></span>
          {{ appState.serviceBusy ? "处理中…" : appViewState.serviceButtonText }}
        </button>
      </div>

      <div v-if="appState.serviceLastError" class="relative mt-5 rounded-[12px] border border-[#ff6b7a]/25 bg-[#ff4455]/10 px-4 py-3 text-sm text-[#ffabb3]">
        {{ appState.serviceLastError }}
      </div>
    </section>

    <section class="mt-4 grid grid-cols-3 gap-3">
      <article class="rounded-[16px] border border-white/8 bg-[#111620] p-4">
        <div class="flex items-center justify-between text-xs text-[#7f8aa3]">
          <span>已配置模型</span>
          <span class="icon-[mdi--cube-outline] text-base text-[#8b7cff]"></span>
        </div>
        <div class="mt-3 text-2xl font-semibold">{{ configuredModels }}</div>
        <div class="mt-1 text-xs text-[#59637a]">自动同步至 Cursor</div>
      </article>
      <article class="rounded-[16px] border border-white/8 bg-[#111620] p-4">
        <div class="flex items-center justify-between text-xs text-[#7f8aa3]">
          <span>对话轮次</span>
          <span class="icon-[mdi--message-processing-outline] text-base text-[#33d6b3]"></span>
        </div>
        <div class="mt-3 text-2xl font-semibold">{{ appState.homeMetrics?.turnsTotal || 0 }}</div>
        <div class="mt-1 text-xs text-[#59637a]">本机使用统计</div>
      </article>
      <article class="rounded-[16px] border border-white/8 bg-[#111620] p-4">
        <div class="flex items-center justify-between text-xs text-[#7f8aa3]">
          <span>Token 使用</span>
          <span class="icon-[mdi--chart-donut] text-base text-[#ffaf65]"></span>
        </div>
        <div class="mt-3 truncate text-2xl font-semibold">{{ tokenTotal }}</div>
        <div class="mt-1 text-xs text-[#59637a]">请求累计</div>
      </article>
    </section>

    <section class="mt-4 rounded-[18px] border border-white/8 bg-[#111620] p-5">
      <div class="flex items-center justify-between gap-4">
        <div class="min-w-0">
          <div class="flex items-center gap-2">
            <span class="flex size-9 items-center justify-center rounded-[10px] bg-[#7568ff]/15 text-[#9186ff]">
              <span class="icon-[mdi--server-network] text-xl"></span>
            </span>
            <div>
              <h2 class="text-sm font-medium text-[#ecf0fb]">模型服务</h2>
              <p class="mt-0.5 max-w-[350px] truncate text-xs text-[#707c96]">{{ activeURL }}</p>
            </div>
          </div>
        </div>
        <button
          class="rounded-[11px] bg-[#7568ff] px-4 py-2.5 text-sm font-medium shadow-[0_8px_22px_rgba(117,104,255,.28)] transition hover:bg-[#8378ff]"
          @click="quickModelSetupVisible = true"
        >
          自动配置模型
        </button>
      </div>

      <div class="my-4 h-px bg-white/7"></div>

      <div class="flex items-center justify-between gap-4">
        <div>
          <div class="text-sm text-[#d9deeb]">高级设置</div>
          <div class="mt-1 text-xs text-[#667189]">查看日志、端口与本地配置文件</div>
        </div>
        <button class="rounded-[10px] border border-white/10 bg-white/5 px-4 py-2 text-sm text-[#b8c0d3] transition hover:bg-white/10" @click="handleOpenConfig">
          打开设置
        </button>
      </div>
    </section>

    <QuickModelSetupModal
      :visible="quickModelSetupVisible"
      @close="quickModelSetupVisible = false"
      @saved="message.success('模型配置已保存')"
    />
  </main>
</template>
