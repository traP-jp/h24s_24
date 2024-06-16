import type { Post } from '@/features/api';

export const getReactions = (post: Post) =>
  post.reactions.map((r) => ({ ...r, clicked: post.my_reactions.includes(r.id) }));
