import { useState, type FormEvent, type ReactNode } from 'react';
import { Plus } from 'lucide-react';
import { AccountEmailStrategy, type CreateGPTAccountRequest, type CreateGPTAccountResponse } from '../proto/orchestrator_account';
import { Button, DashboardDialog, Input, Label, ToolbarIconButton, accountCarrierID, api, errorText, useAsyncActionRunner } from '@byte-v-forge/common-ui';

export function CreateAccountForm({ compact, onDone, onError }: { compact?: boolean; onDone: (message: string) => void; onError: (message: string) => void }) {
  const [open, setOpen] = useState(false);
  const runner = useAsyncActionRunner();
  const [values, setValues] = useState({ email: '', password: '', country_code: 'JP', region: '' });
  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!values.email.trim()) return onError('邮箱不能为空');
    await runner.tryRun('create-account', async () => {
      const payload: CreateGPTAccountRequest = { account_id: '', email: values.email.trim(), password: values.password, country_code: values.country_code.trim(), region: values.region.trim(), email_strategy: AccountEmailStrategy.ACCOUNT_EMAIL_STRATEGY_EXPLICIT };
      const resp = await api<CreateGPTAccountResponse>('/api/gpt/accounts', { method: 'POST', body: JSON.stringify(payload) });
      if (resp.error_message) throw new Error(resp.error_message);
      onDone(`创建账号 已提交: ${accountCarrierID(resp.account) || values.email}`);
      setValues({ email: '', password: '', country_code: values.country_code, region: values.region });
      setOpen(false);
    }, { onError: (err) => onError(errorText(err)) });
  }
  return (
    <>
      {compact ? <ToolbarIconButton label="创建 GPT 账号" tone="primary" icon={<Plus size={15} />} onClick={() => setOpen(true)} /> : <Button size="sm" onClick={() => setOpen(true)}><Plus size={15} /> 添加账号</Button>}
      <DashboardDialog open={open} title="创建 GPT 账号" description="私有模块仅保留显式邮箱创建入口。" size="sm" footer={<Button type="submit" form="create-gpt-account-form" disabled={runner.busy || !values.email.trim()}><Plus size={15} /> {runner.busy ? '提交中' : '创建'}</Button>} onOpenChange={setOpen}>
        <form id="create-gpt-account-form" className="grid gap-3" onSubmit={submit}>
          <Field label="邮箱"><Input value={values.email} onChange={(event) => setValues({ ...values, email: event.target.value })} placeholder="email@example.com" /></Field>
          <Field label="密码"><Input type="password" value={values.password} onChange={(event) => setValues({ ...values, password: event.target.value })} placeholder="可空" /></Field>
          <div className="grid gap-2 sm:grid-cols-2">
            <Field label="国家"><Input value={values.country_code} onChange={(event) => setValues({ ...values, country_code: event.target.value })} placeholder="JP" /></Field>
            <Field label="地区"><Input value={values.region} onChange={(event) => setValues({ ...values, region: event.target.value })} placeholder="可空" /></Field>
          </div>
        </form>
      </DashboardDialog>
    </>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return <label className="grid gap-1"><Label>{label}</Label>{children}</label>;
}
