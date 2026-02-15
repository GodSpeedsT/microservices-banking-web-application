export interface DepositType {
    id: number;
    name: string;
    interestRate: number;
    months: number;
    minAmount: number;
    currency: string;
    description?: string;
    active: boolean;
}

export interface DepositResponse {
    id: number;
    amount: number;
    balance: number;
    currency: string;
    interestRate: number;
    startDate: string;
    maturityDate: string;
    depositType: string;
}

export interface DepositRequest {
    depositTypeId: number;
    amount: number;
    currency: string;
}

export interface UserProfile {
    username: string;
    email?: string;
    roles: string[];
}