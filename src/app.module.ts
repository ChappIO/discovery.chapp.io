import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Agent } from './models/Agent';
import { AppController } from './app.controller';
import {
  POSTGRES_DATABASE,
  POSTGRES_HOST,
  POSTGRES_PASSWORD,
  POSTGRES_USER,
} from './config';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      username: POSTGRES_USER,
      password: POSTGRES_PASSWORD,
      database: POSTGRES_DATABASE,
      host: POSTGRES_HOST,
      synchronize: true,
      entities: [Agent],
    }),
    TypeOrmModule.forFeature([Agent]),
  ],
  controllers: [AppController],
  providers: [],
})
export class AppModule {}
