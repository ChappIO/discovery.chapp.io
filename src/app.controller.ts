import {Body, Controller, Get, HttpException, Param, Post, Req} from "@nestjs/common";
import {Agent} from "./models/Agent";
import {MoreThan, Repository} from "typeorm";
import {InjectRepository} from "@nestjs/typeorm";
import {Request} from "express";
import {IsNotEmpty, IsOptional, registerDecorator, ValidationArguments, ValidationOptions} from "class-validator";
import {ApiResponse, ApiTags} from "@nestjs/swagger";

export function ContainsPrimitives(validationOptions?: ValidationOptions) {
    return function (object: Object, propertyName: string) {
        registerDecorator({
            name: 'containsPrimitives',
            target: object.constructor,
            propertyName: propertyName,
            constraints: [],
            options: validationOptions,
            validator: {
                validate(values: any, args: ValidationArguments) {
                    return Object.values(values).every(value => {
                        return typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean' || value === null;
                    })
                },
            },
        });
    };
}

class RegisterAgent {
    @IsNotEmpty()
    agentId: string;
    @IsOptional()
    @ContainsPrimitives({
        message: "Detailed values may only be primitives"
    })
    details: Record<string, string | number | boolean | null>;
}

@ApiTags('Discovery')
@Controller("agents")
export class AppController {
    constructor(@InjectRepository(Agent) private readonly agents: Repository<Agent>) {
    }

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
            throw new HttpException('The serviceId must be lowercase namespaces separated by dots. A namespace may not start with a number.', 400);
        }
    }

    @Get(":serviceId")
    @ApiResponse({
        status: 200,
        description: 'A list of agents on your network which registered themselves within the last 3 minutes',
    })
    public listAgents(@Req() request: Request, @Param("serviceId") serviceId: string): Promise<Agent[]> {
        AppController.checkServiceId(serviceId);
        const publicAddress = AppController.getPublicIp(request);
        return this.agents.find({
            where: {
                publicAddress,
                serviceId,
                lastSeen: MoreThan(new Date(Date.now() - 3 * 60 * 1000))
            }
        });
    }

    @Post(":serviceId")
    public registerAgent(@Req() request: Request, @Param("serviceId") serviceId: string, @Body() registerAgent: RegisterAgent): Promise<Agent> {
        AppController.checkServiceId(serviceId);
        const agent = new Agent();
        agent.publicAddress = AppController.getPublicIp(request);
        agent.serviceId = serviceId;
        agent.agentId = registerAgent.agentId;
        agent.lastSeen = new Date();
        agent.details = registerAgent.details || {};
        return this.agents.save(agent);
    }
}
