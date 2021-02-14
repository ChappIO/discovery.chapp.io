import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Agent } from './models/Agent';
import { AppController } from './app.controller';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      username: 'postgres',
      password: 'postgres',
      database: 'postgres',
      host: 'localhost',
      synchronize: true,
      entities: [Agent],
      logging: 'all',
    }),
    TypeOrmModule.forFeature([Agent]),
  ],
  controllers: [AppController],
  providers: [],
})
export class AppModule {}
