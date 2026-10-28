<script setup lang="ts">
import ThemeToggle from '@/app/layout/ThemeToggle.vue'
import AuthForm from '@/features/auth/ui/AuthForm.vue'
import AuthPromo from '@/features/auth/ui/AuthPromo.vue'
import AppLogo from '@/shared/ui/AppLogo.vue'

// "Вход" from the Auth boards. A desktop has the promo on the left and the
// form on the right; a phone has the logo, the title and the form.
</script>

<template>
  <div class="auth">
    <!-- A div of its own, so this page's CSS does not reach AuthPromo's root -->
    <div class="promo-slot">
      <AuthPromo />
    </div>

    <section class="side">
      <div class="top">
        <RouterLink to="/" class="brand fs-link" aria-label="FlowScore, на главную">
          <AppLogo :size="26" />
          <span class="wordmark" aria-hidden="true">FLOWSCORE</span>
        </RouterLink>
        <ThemeToggle variant="boxed" />
      </div>
      <div class="column fs-in">
        <h1 class="title">Почувствуй, кто <span class="accent">владеет игрой</span></h1>
        <AuthForm />
      </div>
    </section>
  </div>
</template>

<style scoped>
/* Phone, as on MobileAuth */
.auth {
  min-height: 100dvh;
  background: var(--fs-bg);
}

.promo-slot {
  display: none;
}

.side {
  display: flex;
  flex-direction: column;
  gap: 22px;
  /* The board's 66 px is the status bar and 20 px under it */
  padding: calc(var(--fs-safe-top) + 20px) 20px calc(var(--fs-safe-bottom) + 28px);
}

.top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  gap: 10px;
  align-items: center;
  color: var(--fs-text);
  text-decoration: none;
}

.wordmark {
  font-family: var(--fs-font-display);
  font-size: 22px;
  font-weight: 800;
  letter-spacing: 0.05em;
}

.column {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.title {
  margin-top: var(--fs-space-8);
  font-size: var(--fs-text-large-title);
  font-weight: 700;
  line-height: 1.05;
  letter-spacing: -0.025em;
}

.accent {
  color: var(--fs-home);
}

/* Desktop, as on Auth: the promo, and the form alone in the middle */
@media (min-width: 768px) {
  .auth {
    display: grid;
    grid-template-columns: minmax(0, 1.1fr) minmax(0, 1fr);
  }

  .promo-slot {
    display: flex;
  }

  .side {
    position: relative;
    align-items: center;
    justify-content: center;
    padding: 48px;
  }

  .top {
    position: absolute;
    top: var(--fs-space-24);
    right: var(--fs-space-24);
  }

  .brand,
  .title {
    display: none;
  }

  .column {
    width: 420px;
  }
}
</style>
