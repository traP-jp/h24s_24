import type { Reaction } from '@/features/api';

export const reactionIcons = ['🩷', '🔥', '💧', '😢', '🤔'];

export const convertReactions = (src: Reaction[], my: number[]) => {
  const dist: { id: number; count: number; clicked: boolean }[] = [];
  for (let i = 0; i < reactionIcons.length; i++) {
    const found = src.find((r) => r.id == i);
    if (found) {
      dist.push({
        id: i,
        count: found.count,
        clicked: my.find((m) => m == i) != undefined,
      });
    } else {
      dist.push({
        id: i,
        count: 0,
        clicked: false,
      });
    }
  }
  return dist;
};
