export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

export const fetchApi = async (
  method: HttpMethod,
  path: string,
  option?: { parameters?: Record<string, string | undefined>; body?: Record<string, unknown> },
) => {
  const bodyObj = option?.body && {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(option?.body),
  };
  const parameterStr = option?.parameters
    ? new URLSearchParams(JSON.parse(JSON.stringify(option?.parameters))).toString()
    : '';
  const res = await fetch(`/api${path}?${parameterStr}`, {
    method,
    ...bodyObj,
  });
  const data = await res.json();
  return data;
};

export type Expand<T> = T extends infer O ? { [K in keyof O]: O[K] } : never;

export type CreatedPost = {
  /**
   * 変換元のメッセージ
   */
  original_message: string;
  /**
   * 変換後のメッセージ
   */
  converted_message: string;
  /**
   * UUID
   */
  id: string;
  /**
   * 投稿時刻
   */
  created_at: string;
  /**
   * リプライのとき、親のID。そうでなければ自身のID
   */
  parent_id: string;
  /**
   * リプライのとき、そのおおもとのID。そうでなければ自身のID
   */
  root_id: string;
};

export type Post = Omit<CreatedPost, 'parent_id'> & {
  /**
   * 投稿したユーザー名
   */
  user_name: string;
  /**
   * リアクションのリスト
   */
  reactions: Reaction[];
  /**
   * 自分がリアクションしたリアクションIDのリスト
   */
  my_reactions: number[];
};
export type PostWithoutParents = Omit<Post, 'root_id'>;
export type PostDetail = PostWithoutParents & {
  /**
   * 全ての祖先投稿で、古い順
   */
  ancestors: Array<{
    post: Omit<PostWithoutParents, 'parent_id' | 'root_id'>;
    children_count: number;
  }>;
  /**
   * 1個下の子投稿で、新しい順
   */
  children: Array<{
    post: Omit<PostWithoutParents, 'parent_id' | 'root_id'>;
    children_count: number;
  }>;
};

export type Reaction = {
  /**
   * リアクションID
   */
  id: number;
  /**
   * カウント
   */
  count: number;
};
export type ReactionDetail = {
  /**
   * リアクションID
   */
  id: number;
  /**
   * リアクションしたユーザーのID
   */
  users: string[];
};

export type CreatePostBody = {
  /**
   * メッセージ
   */
  message: string;
  /**
   * リプライのとき、親のID。そうでなければundefined
   */
  parent_id?: string;
};
export type CreatePostResponse = CreatedPost;
export const createPost = async (body: CreatePostBody): Promise<CreatePostResponse> => {
  return fetchApi('POST', '/posts', { body });
};

export type GetPostsParameters = {
  /**
   * 取得件数。デフォルト30
   */
  limit?: number;
  /**
   * このIDの投稿より後に投稿されたものを取得する。指定されない場合は、最新のものからlimit件取得する
   */
  after?: string;
  /**
   * リポストのやつを含むかどうか。デフォルトはfalse
   */
  repost?: boolean;
};
export type GetPostsResponse = Array<
  Expand<
    Post & {
      /**
       * リポストの場合はリポストしたユーザーの名前
       */
      repost_user?: string;
    }
  >
>;
export const getPosts = async ({
  limit,
  after,
  repost,
}: GetPostsParameters): Promise<GetPostsResponse> => {
  return fetchApi('GET', '/posts', {
    parameters: { limit: limit?.toString() ?? '30', after, repost: repost?.toString() ?? 'false' },
  });
};

export type GetPostResponse = Expand<PostDetail>;
export const getPost = async (postId: string): Promise<GetPostResponse> => {
  return fetchApi('GET', `/posts/${postId}`);
};

export type PostReactionResponse = Reaction[];
export const postReaction = async (
  postId: string,
  reactionId: number,
): Promise<PostReactionResponse> => {
  return fetchApi('POST', `/posts/${postId}/reactions/${reactionId}`);
};

export type DeleteReactionResponse = Reaction[];
export const deleteReaction = async (
  postId: string,
  reactionId: number,
): Promise<DeleteReactionResponse> => {
  return fetchApi('DELETE', `/posts/${postId}/reactions/${reactionId}`);
};

export type GetReactionsResponse = ReactionDetail[];
export const getReactions = async (postId: string): Promise<GetReactionsResponse> => {
  return fetchApi('GET', `/posts/${postId}/reactions`);
};

export type GetTrendResponse = Array<Post>;
export const getTrend = async (reactionId: number): Promise<GetTrendResponse> => {
  return fetchApi('GET', '/trend', { parameters: { reaction_id: reactionId.toString() } });
};

export type GetUserResponse = {
  /**
   * ユーザーID
   */
  user_name: string;
  /**
   * 投稿数
   */
  post_count: number;
  /**
   * リアクションした数
   */
  reaction_count: number;
  /**
   * リアクションされた数
   */
  get_reaction_count: number;
  /**
   * 投稿のリスト
   */
  posts: Post[];
};
export const getUser = async (userName: string): Promise<GetUserResponse> => {
  return fetchApi('GET', `/user/${userName}`);
};
