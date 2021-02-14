import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { NestExpressApplication } from '@nestjs/platform-express';

async function bootstrap() {
  const app = await NestFactory.create<NestExpressApplication>(AppModule, {
    logger: true,
  });
  app.disable('x-powered-by');
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      transform: true,
    }),
  );

  const document = SwaggerModule.createDocument(
    app,
    new DocumentBuilder()
      .setTitle('Agent discovery service')
      .setDescription('A simple way to discover agents on your local network')
      .setVersion('2.0')
      .addTag('Discovery')
      .build(),
  );
  SwaggerModule.setup('', app, document);
  await app.listen(3000);
}

bootstrap().catch((e) => {
  console.error(e);
  process.exit(1);
});
