import {
  Body,
  Controller,
  Get,
  HttpException,
  Param,
  Post,
  Req,
} from '@nestjs/common';
import { Agent } from './models/Agent';
import { MoreThan, Repository } from 'typeorm';
import { InjectRepository } from '@nestjs/typeorm';
import { Request } from 'express';
import {
  IsNotEmpty,
  IsOptional,
  registerDecorator,
  ValidationOptions,
} from 'class-validator';
import {
  ApiBody,
  ApiOperation,
  ApiParam,
  ApiProperty,
  ApiPropertyOptional,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';

export function ContainsPrimitives(validationOptions?: ValidationOptions) {
  return function (object: any, propertyName: string) {
    registerDecorator({
      name: 'containsPrimitives',
      target: object.constructor,
      propertyName: propertyName,
      constraints: [],
      options: validationOptions,
      validator: {
        validate(values: any) {
          return Object.values(values).every((value) => {
            return (
              typeof value === 'string' ||
              typeof value === 'number' ||
              typeof value === 'boolean' ||
              value === null
            );
          });
        },
      },
    });
  };
}

class RegisterAgent {
  @ApiProperty({
    description: 'A unique identifier for your agent on this network.',
    example: 'b3dd59db-0fdc-4a2d-a730-80a861568021',
  })
  @IsNotEmpty()
  agentId: string;
  @IsOptional()
  @ContainsPrimitives({
    message: 'Detailed values may only be primitives',
  })
  @ApiPropertyOptional({
    description:
      'A free-form object of metadata to attach to this agent. The values of this object may only be primitives.',
    example: {
      displayName: "Thomas' Desktop",
      customField: 'Custom Value',
    },
  })
  details: Record<string, string | number | boolean | null>;
}

@ApiTags('Discovery')
@Controller('agents')
export class AppController {
  constructor(
    @InjectRepository(Agent) private readonly agents: Repository<Agent>,
  ) {}

  private static getPublicIp(request: Request): string {
    const forwardedFor = request.header('x-forwarded-for');
    if (forwardedFor) {
      const firstIp = forwardedFor.split(',')[0];
      if (firstIp) {
        return firstIp.trim();
      }
    }
    return request.socket.remoteAddress;
  }

  private static checkServiceId(serviceId: string) {
    if (!serviceId.match(/^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$/)) {
      throw new HttpException(
        'The serviceId must be lowercase namespaces separated by dots. A namespace may not start with a number.',
        400,
      );
    }
  }

  @ApiOperation({
    summary: 'Register as an agent',
    description:
      'Register or update an agent. This operation is idempotent and designed such that your agent can call this endpoint every two minutes.',
  })
  @ApiBody({
    type: RegisterAgent,
  })
  @ApiParam({
    name: 'serviceId',
    description:
      'The unique identifier of your service. Insert your development namespace here. ' +
      'The namespace must consist of one or more parts separated by dots. Each part must be a lowercase alphanumeric string and may not start with a number.',
    example: 'io.chapp.myproject',
  })
  @Post(':serviceId')
  public registerAgent(
    @Req() request: Request,
    @Param('serviceId') serviceId: string,
    @Body() registerAgent: RegisterAgent,
  ): Promise<Agent> {
    AppController.checkServiceId(serviceId);
    const agent = new Agent();
    agent.publicAddress = AppController.getPublicIp(request);
    agent.serviceId = serviceId;
    agent.agentId = registerAgent.agentId;
    agent.lastSeen = new Date();
    agent.details = registerAgent.details || {};
    return this.agents.save(agent);
  }

  @Get(':serviceId')
  @ApiResponse({
    status: 200,
    description:
      'A list of agents on your network which registered themselves within the last 3 minutes',
  })
  @ApiParam({
    name: 'serviceId',
    description:
      'The unique identifier of your service. Insert your development namespace here. ' +
      'The namespace must consist of one or more parts separated by dots. Each part must be a lowercase alphanumeric string and may not start with a number.',
    example: 'io.chapp.myproject',
  })
  public listAgents(
    @Req() request: Request,
    @Param('serviceId') serviceId: string,
  ): Promise<Agent[]> {
    AppController.checkServiceId(serviceId);
    const publicAddress = AppController.getPublicIp(request);
    return this.agents.find({
      where: {
        publicAddress,
        serviceId,
        lastSeen: MoreThan(new Date(Date.now() - 3 * 60 * 1000)),
      },
    });
  }
}
