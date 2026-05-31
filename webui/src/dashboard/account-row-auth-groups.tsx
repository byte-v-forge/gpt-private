import { Play } from 'lucide-react';
import { AccountRowActionGroups, accountActionAvailability, accountActionLabel } from '@byte-v-forge/common-ui';
import type { RowActionDescriptor } from '@byte-v-forge/common-ui';
import { GPT_ACTIONS, type GptActionCatalog } from './action-catalog';
import { canRegister } from './account-utils';
import type { Account } from './types';

export function AccountRowAuthGroups({ account, actionCatalog, busy, onRegisterProtocol }: {
  account: Account;
  actionCatalog?: GptActionCatalog;
  busy: boolean;
  onRegisterProtocol: (account: Account) => void;
}) {
  const actions = rowActions(account, actionCatalog, busy, onRegisterProtocol);
  return <AccountRowActionGroups actions={actions} />;
}

function rowActions(account: Account, catalog: GptActionCatalog | undefined, busy: boolean, onRegisterProtocol: (account: Account) => void): RowActionDescriptor[] {
  const availability = accountActionAvailability(catalog, GPT_ACTIONS.registerProtocol, account, 'account_row');
  if (!availability.visible) return [];
  return [{
    label: accountActionLabel(catalog, GPT_ACTIONS.registerProtocol, '注册', 'account_row'),
    icon: <Play size={14} />,
    onClick: () => onRegisterProtocol(account),
    disabled: busy || !availability.enabled || !canRegister(account),
    kind: 'secondary'
  }];
}
