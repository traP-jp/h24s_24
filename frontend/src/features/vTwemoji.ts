import twemoji from 'twemoji';

export const vTwemojiObj = {
  mounted: (el: HTMLElement) => {
    el.innerHTML = twemoji.parse(el.innerHTML, {
      className: 'twemoji',
    });
  },
};
