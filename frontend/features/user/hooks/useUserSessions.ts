import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getUserSessions,
  deleteUserSession,
  getInactivityTimeout,
  setInactivityTimeout,
} from '@/features/user/api/user';

export const useUserSessions = () =>
  useQuery({ queryKey: ['sessions'], queryFn: getUserSessions });

export const useDeleteSession = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteUserSession,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['sessions'] }),
  });
};

export const useInactivityTimeout = () =>
  useQuery({ queryKey: ['inactivity-timeout'], queryFn: getInactivityTimeout });

export const useSetInactivityTimeout = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: setInactivityTimeout,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['inactivity-timeout'] }),
  });
};
