import { ref, type UnwrapRef } from 'vue';

export const useFetcher = <T>(fetcher: () => Promise<UnwrapRef<T>>) => {
  const data = ref<T | undefined>(undefined);
  const loading = ref(true);
  const error = ref<Error | undefined>(undefined);

  fetcher()
    .then((fetchedData) => {
      data.value = fetchedData;
      loading.value = false;
    })
    .catch((e) => {
      error.value = e;
      loading.value = false;
    });

  return { data, loading, error };
};
