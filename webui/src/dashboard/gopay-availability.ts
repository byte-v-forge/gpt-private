import { api, useQuery } from '@byte-v-forge/common-ui';

export type GoPayAvailability = {
  available: boolean;
  loading: boolean;
  reason: string;
};

type GoPayHealthResponse = {
  success?: boolean;
  ok?: boolean;
  service?: string;
};

export function useGoPayAvailability(): GoPayAvailability {
  const query = useQuery({
    queryKey: ['gopay', 'availability'],
    queryFn: () => api<GoPayHealthResponse>('/api/gopay/health'),
    retry: false,
    staleTime: 15000,
    refetchInterval: 30000
  });
  const available = Boolean(query.data?.success || query.data?.ok);
  return {
    available,
    loading: query.isLoading,
    reason: available ? '' : (query.isLoading ? '检查 gopay-app...' : 'gopay-app 不可用')
  };
}
