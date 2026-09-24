# FlowScore

Live-аналитика футбольных матчей: метрика Flow Momentum плюс AI-разбор.

- Архитектура и жизненный цикл: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- Правила фронтенда: [FRONTEND.md](FRONTEND.md). **Обязательны для любой задачи по фронту.**

## Коротко о фронте
- Мокап в `design/` соблюдается на 100%. Токены только через `--fs-*`. У каждого экрана 5 состояний.
- Главный скилл `apple-design`, остальные фронтовые скиллы вторичны.
- Lenis (только десктоп, на тикере GSAP, синхронизирован со ScrollTrigger), GSAP + ScrollTrigger, Rive, Three.js (лениво).
- Easing `cubic-bezier(0.2, 0.7, 0.2, 1)`, длительности 150/300/700 мс, поддержка reduced-motion.
