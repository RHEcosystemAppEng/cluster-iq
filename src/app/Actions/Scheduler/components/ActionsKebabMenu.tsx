import { api, ActionResponseApi } from '@api';

export const rowActions = (action: ActionResponseApi, reloadActions: () => Promise<void>) => [
  {
    title: 'Enable',
    onClick: async () => {
      await api.actions.actionsEnable(action.id);
      await reloadActions();
    },
  },
  {
    title: 'Disable',
    onClick: async () => {
      await api.actions.actionsDisable(action.id);
      await reloadActions();
    },
  },
  { isSeparator: true },
  {
    title: 'Delete',
    isDanger: true,
    onClick: async () => {
      await api.actions.actionsDelete(action.id);
      await reloadActions();
    },
  },
];
