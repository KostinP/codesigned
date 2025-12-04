export interface UserSession {
  id: string;
  user_id: string;
  token: string;
  ip_address: string;
  city: string;
  country: string;
  user_agent: string;
  created_at: string;
  expires_at: string;
  last_active_at: string;
}
