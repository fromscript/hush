FROM node:22-alpine AS deps
WORKDIR /app

RUN apk add --no-cache libc6-compat

# Copy only the files needed to install dependencies
COPY package.json package-lock.json* ./

# Install dependencies with the preferred package manager
RUN npm ci


FROM node:22-alpine AS builder
WORKDIR /app

ARG NEXT_PUBLIC_WS_URL
ENV NEXT_PUBLIC_WS_URL=$NEXT_PUBLIC_WS_URL

COPY --from=deps /app/node_modules ./node_modules

# Copy the rest of the files
COPY . .

# Run build with the preferred package manager
RUN npm run build


FROM node:22-alpine AS runner
WORKDIR /app

ENV NODE_ENV production

# Add nextjs user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy public assets
#COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next
RUN chown nextjs:nodejs .next

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

ENV PORT 8080
EXPOSE 8080
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]
