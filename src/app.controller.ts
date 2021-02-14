import {Body, Controller, Get, HttpException, Param, Post, Req} from "@nestjs/common";
import {Agent} from "./models/Agent";
import {Repository} from "typeorm";
import {InjectRepository} from "@nestjs/typeorm";
import {Request} from "express";
import {
    IsEmail,
    IsNotEmpty,
    IsOptional,
    registerDecorator,
    ValidationArguments,
    ValidationOptions
} from "class-validator";

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
        if(!serviceId.match(/^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$/)) {
            throw new HttpException('The serviceId must be lowercase namespaces separated by dots. A namespace may not start with a number.', 400);
        }
    }

    @Get(":serviceId")
    public listAgents(@Req() request: Request, @Param("serviceId") serviceId: string): Promise<Agent[]> {
        AppController.checkServiceId(serviceId);
        const publicAddress = AppController.getPublicIp(request);
        return this.agents.find({
            publicAddress,
            serviceId
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

        console.log(registerAgent);

        return this.agents.save(agent);
    }
}
