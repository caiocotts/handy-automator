import AsyncStorage from "@react-native-async-storage/async-storage";

const BASE_URL = process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080/api";

const ACCESS_TOKEN_KEY = "accessToken";
const REFRESH_TOKEN_KEY = "refreshToken";

export async function getAccessToken(): Promise<string | null> {
    return AsyncStorage.getItem(ACCESS_TOKEN_KEY);
}

export async function saveTokens(
    accessToken: string,
    refreshToken: string,
): Promise<void> {
    await AsyncStorage.multiSet([
        [ACCESS_TOKEN_KEY, accessToken],
        [REFRESH_TOKEN_KEY, refreshToken],
    ]);
}

export async function clearTokens(): Promise<void> {
    await AsyncStorage.multiRemove([ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY]);
}

async function authFetch(
    path: string,
    options: RequestInit = {},
): Promise<Response> {
    const token = await getAccessToken();
    const headers: Record<string, string> = {
        "Content-Type": "application/json",
        ...(options.headers as Record<string, string>),
    };
    if (token) {
        headers["Authorization"] = `Bearer ${token}`;
    }
    return fetch(`${BASE_URL}${path}`, { ...options, headers });
}

// Auth

export interface LoginResponse {
    userId: string;
    username: string;
    accessToken: string;
    refreshToken: string;
}

export async function login(
    username: string,
    password: string,
): Promise<LoginResponse> {
    const res = await fetch(`${BASE_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    });
    if (res.status === 401) throw new Error("Invalid credentials");
    if (!res.ok) throw new Error("Login failed");
    return res.json();
}

// Workflows

export interface Device {
    id: string;
    hostname: string;
    name?: string;
    type?: string;
}

export interface Workflow {
    id: string;
    name: string;
    userId: string;
    devices?: Device[];
}

export interface DeviceTriggerStatus {
    deviceId: string;
    ok: boolean;
    error?: string;
}

export async function getWorkflows(): Promise<Workflow[]> {
    const res = await authFetch("/workflow");
    if (!res.ok) throw new Error("Failed to fetch workflows");
    const data = await res.json();
    return data.workflows ?? [];
}

export async function getWorkflow(id: string): Promise<Workflow> {
    const res = await authFetch(`/workflow/${id}`);
    if (!res.ok) throw new Error("Failed to fetch workflow");
    return res.json();
}

export async function createWorkflow(name: string): Promise<Workflow> {
    const res = await authFetch("/workflow", {
        method: "POST",
        body: JSON.stringify({ name }),
    });
    if (!res.ok) throw new Error("Failed to create workflow");
    return res.json();
}

export async function deleteWorkflow(id: string): Promise<void> {
    const res = await authFetch(`/workflow/${id}`, { method: "DELETE" });
    if (!res.ok) throw new Error("Failed to delete workflow");
}

export async function triggerWorkflow(
    id: string,
): Promise<DeviceTriggerStatus[]> {
    const res = await authFetch(`/workflow/${id}/trigger`, { method: "POST" });
    if (!res.ok) throw new Error("Failed to trigger workflow");
    return res.json();
}

export async function getDevices(): Promise<Device[]> {
    const res = await authFetch("/device");
    if (!res.ok) throw new Error("Failed to fetch devices");
    const data = await res.json();
    return data.devices ?? [];
}

export async function associateWorkflowDevices(
    workflowId: string,
    deviceIds: string[],
): Promise<string[]> {
    const res = await authFetch(`/workflow/${workflowId}/devices`, {
        method: "PUT",
        body: JSON.stringify({ devices: deviceIds }),
    });
    if (!res.ok) throw new Error("Failed to associate devices");
    const data = await res.json();
    return data.devices ?? [];
}
