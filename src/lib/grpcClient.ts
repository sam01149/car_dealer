// Impor service client Anda
import { AuthServiceClient } from "@/proto/CarappServiceClientPb";
import { NhtsaDataServiceClient } from "@/proto/CarappServiceClientPb";
import { MobilServiceClient } from "@/proto/CarappServiceClientPb";
import { TransaksiServiceClient } from "@/proto/CarappServiceClientPb";
import { NotifikasiServiceClient } from "@/proto/CarappServiceClientPb";
import { DashboardServiceClient } from "@/proto/CarappServiceClientPb";

const API_URL = "http://localhost:9090";

// --- INTERCEPTOR (Pencegat Otomatis) ---
// Fungsi helper untuk menambahkan token ke metadata
function addAuthMetadata(metadata: any = {}) {
  // Cek apakah ada token di localStorage
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('authToken');
    console.log('🔑 addAuthMetadata called');
    console.log('🔑 Token from localStorage:', token ? `${token.substring(0, 20)}...` : 'NULL');
    
    if (token) {
      // Untuk grpc-web, header harus lowercase
      metadata['authorization'] = `Bearer ${token}`;
      console.log('✅ Token added to metadata');
    } else {
      console.warn('⚠️ No token found in localStorage!');
    }
  }
  return metadata;
}

// Wrapper untuk client dengan auto-inject token
class AuthInterceptor {
  intercept(request: any, metadata: any) {
    return addAuthMetadata(metadata);
  }
}

const interceptor = new AuthInterceptor();

// Options untuk grpc-web client
const options = {};

// Kita ekspor instance klien untuk setiap service
export const authClient = new AuthServiceClient(API_URL, null, options);
export const nhtsaClient = new NhtsaDataServiceClient(API_URL, null, options);
export const mobilClient = new MobilServiceClient(API_URL, null, options);
export const transaksiClient = new TransaksiServiceClient(API_URL, null, options);
export const notifikasiClient = new NotifikasiServiceClient(API_URL, null, options);
export const dashboardClient = new DashboardServiceClient(API_URL, null, options);

// Export helper untuk menambahkan auth metadata
export { addAuthMetadata };

