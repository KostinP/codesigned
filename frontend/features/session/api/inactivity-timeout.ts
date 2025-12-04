import { axios } from '@/lib/axios' 

export const getInactivityTimeout = async (): Promise<{ timeout: number }> => {
  const { data } = await axios.get('/me/inactivity-timeout')
  return data
}

export const setInactivityTimeout = async (timeoutSeconds: number) => {
  await axios.post('/me/inactivity-timeout', {
    timeout_seconds: timeoutSeconds,
  })
}
