FROM node:22-alpine AS base

FROM base AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
ARG NEXT_PUBLIC_SITE_URL=http://localhost:3000
ARG NEXT_PUBLIC_SITE_NAME="Под Покровом Божией Матери \"Тихвинская\""
ARG NEXT_PUBLIC_SITE_FULL_NAME=РПЦ Тихвинская Епархия. Паломническая служба «Под покровом Божией Матери «Тихвинская»»
ARG NEXT_PUBLIC_SITE_TAGLINE=Паломническая служба
ARG NEXT_PUBLIC_SITE_DESCRIPTION=Паломнические поездки в монастыри и святые места России. Сопровождение духовника, комфортный транспорт, прозрачная стоимость.
ARG NEXT_PUBLIC_CONTACT_PHONE=+79669334321
ARG NEXT_PUBLIC_CONTACT_PHONE_DISPLAY=+7 966 933-43-21
ARG NEXT_PUBLIC_CONTACT_EMAIL=info@tikhvin-palomnik.ru
ARG NEXT_PUBLIC_YM_ID=
ARG NEXT_PUBLIC_GA_ID=
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_SITE_URL=$NEXT_PUBLIC_SITE_URL
ENV NEXT_PUBLIC_SITE_NAME=$NEXT_PUBLIC_SITE_NAME
ENV NEXT_PUBLIC_SITE_FULL_NAME=$NEXT_PUBLIC_SITE_FULL_NAME
ENV NEXT_PUBLIC_SITE_TAGLINE=$NEXT_PUBLIC_SITE_TAGLINE
ENV NEXT_PUBLIC_SITE_DESCRIPTION=$NEXT_PUBLIC_SITE_DESCRIPTION
ENV NEXT_PUBLIC_CONTACT_PHONE=$NEXT_PUBLIC_CONTACT_PHONE
ENV NEXT_PUBLIC_CONTACT_PHONE_DISPLAY=$NEXT_PUBLIC_CONTACT_PHONE_DISPLAY
ENV NEXT_PUBLIC_CONTACT_EMAIL=$NEXT_PUBLIC_CONTACT_EMAIL
ENV NEXT_PUBLIC_YM_ID=$NEXT_PUBLIC_YM_ID
ENV NEXT_PUBLIC_GA_ID=$NEXT_PUBLIC_GA_ID
ENV NEXT_TELEMETRY_DISABLED=1

RUN npm run build

FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000

CMD ["node", "server.js"]
