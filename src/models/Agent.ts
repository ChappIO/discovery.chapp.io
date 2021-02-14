import {Column, Entity, PrimaryColumn, UpdateDateColumn} from "typeorm";

export type AgentParameters = Record<string, string | number | boolean | null>;

@Entity()
export class Agent {
    @PrimaryColumn()
    serviceId: string;
    @PrimaryColumn()
    publicAddress: string;
    @PrimaryColumn()
    agentId: string;
    @Column()
    lastSeen: Date;
    @Column({
        type: 'json',
    })
    details: AgentParameters;
}
