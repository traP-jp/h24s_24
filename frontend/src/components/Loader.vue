<template>
  <div class="loader"></div>
</template>

<style scoped>
.loader {
  width: 32px;
  aspect-ratio: 1;
  --c: no-repeat linear-gradient(var(--accent-color) 0 0);
  background:
    var(--c) 0 0,
    var(--c) 50% 50%,
    var(--c) 100% 100%;
  animation:
    l17 1s infinite alternate;
}

@keyframes l17 {

  0%,
  10% {
    background-size: 20% 100%
  }

  50% {
    background-size: 20% 20%
  }

  90%,
  100% {
    background-size: 100% 20%
  }
}
</style>
