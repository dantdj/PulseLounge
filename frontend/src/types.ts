export type Message = {
  id: number;
  author_id: number;
  channel_id: number;
  body: string;
  image?: string;
  created_at: string;
  edited_at: string | null;
};

export type Channel = {
  id: number;
  name: string;
  created_at: string;
};
