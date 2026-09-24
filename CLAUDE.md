# FlowScore

Live-аналитика футбольных матчей: метрика Flow Momentum плюс AI-разбор.

- Архитектура, источники данных и жизненный цикл: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- Формула Flow и модель xG: [docs/FLOW.md](docs/FLOW.md)
- Правила фронтенда: [FRONTEND.md](FRONTEND.md). **Обязательны для любой задачи по фронту.**

## Коротко о фронте
- Мокап в `design/` соблюдается на 100%. Токены только через `--fs-*`. У каждого экрана 5 состояний.
- Главный скилл `apple-design`, остальные фронтовые скиллы вторичны.
- Lenis (только десктоп, на тикере GSAP, синхронизирован со ScrollTrigger), GSAP + ScrollTrigger, Rive, Three.js (лениво).
- Easing `cubic-bezier(0.2, 0.7, 0.2, 1)`, длительности только из токенов `--fs-duration-*` (таблица в FRONTEND.md), поддержка reduced-motion.
- Мокап лежит в `design/mockup/`, сверяться с ним.
