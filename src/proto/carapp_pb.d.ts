import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"


export class User extends jspb.Message {
  getId(): string;
  setId(value: string): User;

  getName(): string;
  setName(value: string): User;

  getEmail(): string;
  setEmail(value: string): User;

  getPhone(): string;
  setPhone(value: string): User;

  getRole(): string;
  setRole(value: string): User;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): User;
  hasCreatedAt(): boolean;
  clearCreatedAt(): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: string,
    name: string,
    email: string,
    phone: string,
    role: string,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class Mobil extends jspb.Message {
  getId(): string;
  setId(value: string): Mobil;

  getOwnerId(): string;
  setOwnerId(value: string): Mobil;

  getMerk(): string;
  setMerk(value: string): Mobil;

  getModel(): string;
  setModel(value: string): Mobil;

  getTahun(): number;
  setTahun(value: number): Mobil;

  getKondisi(): string;
  setKondisi(value: string): Mobil;

  getDeskripsi(): string;
  setDeskripsi(value: string): Mobil;

  getHargaJual(): number;
  setHargaJual(value: number): Mobil;

  getFotoUrl(): string;
  setFotoUrl(value: string): Mobil;

  getLokasi(): string;
  setLokasi(value: string): Mobil;

  getStatus(): string;
  setStatus(value: string): Mobil;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Mobil;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Mobil;

  getOwnerName(): string;
  setOwnerName(value: string): Mobil;

  getHargaRentalPerHari(): number;
  setHargaRentalPerHari(value: number): Mobil;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Mobil.AsObject;
  static toObject(includeInstance: boolean, msg: Mobil): Mobil.AsObject;
  static serializeBinaryToWriter(message: Mobil, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Mobil;
  static deserializeBinaryFromReader(message: Mobil, reader: jspb.BinaryReader): Mobil;
}

export namespace Mobil {
  export type AsObject = {
    id: string,
    ownerId: string,
    merk: string,
    model: string,
    tahun: number,
    kondisi: string,
    deskripsi: string,
    hargaJual: number,
    fotoUrl: string,
    lokasi: string,
    status: string,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    ownerName: string,
    hargaRentalPerHari: number,
  }
}

export class Notifikasi extends jspb.Message {
  getId(): string;
  setId(value: string): Notifikasi;

  getUserId(): string;
  setUserId(value: string): Notifikasi;

  getTipe(): string;
  setTipe(value: string): Notifikasi;

  getPesan(): string;
  setPesan(value: string): Notifikasi;

  getPriority(): string;
  setPriority(value: string): Notifikasi;

  getReadAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setReadAt(value?: google_protobuf_timestamp_pb.Timestamp): Notifikasi;
  hasReadAt(): boolean;
  clearReadAt(): Notifikasi;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Notifikasi;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Notifikasi;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Notifikasi.AsObject;
  static toObject(includeInstance: boolean, msg: Notifikasi): Notifikasi.AsObject;
  static serializeBinaryToWriter(message: Notifikasi, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Notifikasi;
  static deserializeBinaryFromReader(message: Notifikasi, reader: jspb.BinaryReader): Notifikasi;
}

export namespace Notifikasi {
  export type AsObject = {
    id: string,
    userId: string,
    tipe: string,
    pesan: string,
    priority: string,
    readAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class RegisterRequest extends jspb.Message {
  getName(): string;
  setName(value: string): RegisterRequest;

  getEmail(): string;
  setEmail(value: string): RegisterRequest;

  getPassword(): string;
  setPassword(value: string): RegisterRequest;

  getPhone(): string;
  setPhone(value: string): RegisterRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterRequest): RegisterRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterRequest;
  static deserializeBinaryFromReader(message: RegisterRequest, reader: jspb.BinaryReader): RegisterRequest;
}

export namespace RegisterRequest {
  export type AsObject = {
    name: string,
    email: string,
    password: string,
    phone: string,
  }
}

export class LoginRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): LoginRequest;

  getPassword(): string;
  setPassword(value: string): LoginRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LoginRequest): LoginRequest.AsObject;
  static serializeBinaryToWriter(message: LoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginRequest;
  static deserializeBinaryFromReader(message: LoginRequest, reader: jspb.BinaryReader): LoginRequest;
}

export namespace LoginRequest {
  export type AsObject = {
    email: string,
    password: string,
  }
}

export class AuthResponse extends jspb.Message {
  getUser(): User | undefined;
  setUser(value?: User): AuthResponse;
  hasUser(): boolean;
  clearUser(): AuthResponse;

  getToken(): string;
  setToken(value: string): AuthResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AuthResponse): AuthResponse.AsObject;
  static serializeBinaryToWriter(message: AuthResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthResponse;
  static deserializeBinaryFromReader(message: AuthResponse, reader: jspb.BinaryReader): AuthResponse;
}

export namespace AuthResponse {
  export type AsObject = {
    user?: User.AsObject,
    token: string,
  }
}

export class CreateMobilRequest extends jspb.Message {
  getMerk(): string;
  setMerk(value: string): CreateMobilRequest;

  getModel(): string;
  setModel(value: string): CreateMobilRequest;

  getTahun(): number;
  setTahun(value: number): CreateMobilRequest;

  getKondisi(): string;
  setKondisi(value: string): CreateMobilRequest;

  getDeskripsi(): string;
  setDeskripsi(value: string): CreateMobilRequest;

  getHargaJual(): number;
  setHargaJual(value: number): CreateMobilRequest;

  getFotoUrl(): string;
  setFotoUrl(value: string): CreateMobilRequest;

  getLokasi(): string;
  setLokasi(value: string): CreateMobilRequest;

  getHargaRentalPerHari(): number;
  setHargaRentalPerHari(value: number): CreateMobilRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateMobilRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateMobilRequest): CreateMobilRequest.AsObject;
  static serializeBinaryToWriter(message: CreateMobilRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateMobilRequest;
  static deserializeBinaryFromReader(message: CreateMobilRequest, reader: jspb.BinaryReader): CreateMobilRequest;
}

export namespace CreateMobilRequest {
  export type AsObject = {
    merk: string,
    model: string,
    tahun: number,
    kondisi: string,
    deskripsi: string,
    hargaJual: number,
    fotoUrl: string,
    lokasi: string,
    hargaRentalPerHari: number,
  }
}

export class ListMobilRequest extends jspb.Message {
  getPage(): number;
  setPage(value: number): ListMobilRequest;

  getLimit(): number;
  setLimit(value: number): ListMobilRequest;

  getFilterStatus(): string;
  setFilterStatus(value: string): ListMobilRequest;
  hasFilterStatus(): boolean;
  clearFilterStatus(): ListMobilRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMobilRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMobilRequest): ListMobilRequest.AsObject;
  static serializeBinaryToWriter(message: ListMobilRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMobilRequest;
  static deserializeBinaryFromReader(message: ListMobilRequest, reader: jspb.BinaryReader): ListMobilRequest;
}

export namespace ListMobilRequest {
  export type AsObject = {
    page: number,
    limit: number,
    filterStatus?: string,
  }

  export enum FilterStatusCase { 
    _FILTER_STATUS_NOT_SET = 0,
    FILTER_STATUS = 3,
  }
}

export class ListMobilResponse extends jspb.Message {
  getMobilsList(): Array<Mobil>;
  setMobilsList(value: Array<Mobil>): ListMobilResponse;
  clearMobilsList(): ListMobilResponse;
  addMobils(value?: Mobil, index?: number): Mobil;

  getTotal(): number;
  setTotal(value: number): ListMobilResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMobilResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMobilResponse): ListMobilResponse.AsObject;
  static serializeBinaryToWriter(message: ListMobilResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMobilResponse;
  static deserializeBinaryFromReader(message: ListMobilResponse, reader: jspb.BinaryReader): ListMobilResponse;
}

export namespace ListMobilResponse {
  export type AsObject = {
    mobilsList: Array<Mobil.AsObject>,
    total: number,
  }
}

export class GetMobilRequest extends jspb.Message {
  getMobilId(): string;
  setMobilId(value: string): GetMobilRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMobilRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMobilRequest): GetMobilRequest.AsObject;
  static serializeBinaryToWriter(message: GetMobilRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMobilRequest;
  static deserializeBinaryFromReader(message: GetMobilRequest, reader: jspb.BinaryReader): GetMobilRequest;
}

export namespace GetMobilRequest {
  export type AsObject = {
    mobilId: string,
  }
}

export class Make extends jspb.Message {
  getBrandId(): string;
  setBrandId(value: string): Make;

  getName(): string;
  setName(value: string): Make;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Make.AsObject;
  static toObject(includeInstance: boolean, msg: Make): Make.AsObject;
  static serializeBinaryToWriter(message: Make, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Make;
  static deserializeBinaryFromReader(message: Make, reader: jspb.BinaryReader): Make;
}

export namespace Make {
  export type AsObject = {
    brandId: string,
    name: string,
  }
}

export class Model extends jspb.Message {
  getModelId(): string;
  setModelId(value: string): Model;

  getBrandId(): string;
  setBrandId(value: string): Model;

  getName(): string;
  setName(value: string): Model;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Model.AsObject;
  static toObject(includeInstance: boolean, msg: Model): Model.AsObject;
  static serializeBinaryToWriter(message: Model, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Model;
  static deserializeBinaryFromReader(message: Model, reader: jspb.BinaryReader): Model;
}

export namespace Model {
  export type AsObject = {
    modelId: string,
    brandId: string,
    name: string,
  }
}

export class GetMakesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMakesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMakesRequest): GetMakesRequest.AsObject;
  static serializeBinaryToWriter(message: GetMakesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMakesRequest;
  static deserializeBinaryFromReader(message: GetMakesRequest, reader: jspb.BinaryReader): GetMakesRequest;
}

export namespace GetMakesRequest {
  export type AsObject = {
  }
}

export class GetMakesResponse extends jspb.Message {
  getMakesList(): Array<Make>;
  setMakesList(value: Array<Make>): GetMakesResponse;
  clearMakesList(): GetMakesResponse;
  addMakes(value?: Make, index?: number): Make;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMakesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMakesResponse): GetMakesResponse.AsObject;
  static serializeBinaryToWriter(message: GetMakesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMakesResponse;
  static deserializeBinaryFromReader(message: GetMakesResponse, reader: jspb.BinaryReader): GetMakesResponse;
}

export namespace GetMakesResponse {
  export type AsObject = {
    makesList: Array<Make.AsObject>,
  }
}

export class GetModelsForMakeRequest extends jspb.Message {
  getBrandId(): string;
  setBrandId(value: string): GetModelsForMakeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetModelsForMakeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetModelsForMakeRequest): GetModelsForMakeRequest.AsObject;
  static serializeBinaryToWriter(message: GetModelsForMakeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetModelsForMakeRequest;
  static deserializeBinaryFromReader(message: GetModelsForMakeRequest, reader: jspb.BinaryReader): GetModelsForMakeRequest;
}

export namespace GetModelsForMakeRequest {
  export type AsObject = {
    brandId: string,
  }
}

export class GetModelsForMakeResponse extends jspb.Message {
  getModelsList(): Array<Model>;
  setModelsList(value: Array<Model>): GetModelsForMakeResponse;
  clearModelsList(): GetModelsForMakeResponse;
  addModels(value?: Model, index?: number): Model;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetModelsForMakeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetModelsForMakeResponse): GetModelsForMakeResponse.AsObject;
  static serializeBinaryToWriter(message: GetModelsForMakeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetModelsForMakeResponse;
  static deserializeBinaryFromReader(message: GetModelsForMakeResponse, reader: jspb.BinaryReader): GetModelsForMakeResponse;
}

export namespace GetModelsForMakeResponse {
  export type AsObject = {
    modelsList: Array<Model.AsObject>,
  }
}

export class BuyMobilRequest extends jspb.Message {
  getMobilId(): string;
  setMobilId(value: string): BuyMobilRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BuyMobilRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BuyMobilRequest): BuyMobilRequest.AsObject;
  static serializeBinaryToWriter(message: BuyMobilRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BuyMobilRequest;
  static deserializeBinaryFromReader(message: BuyMobilRequest, reader: jspb.BinaryReader): BuyMobilRequest;
}

export namespace BuyMobilRequest {
  export type AsObject = {
    mobilId: string,
  }
}

export class TransaksiJualResponse extends jspb.Message {
  getId(): string;
  setId(value: string): TransaksiJualResponse;

  getMobilId(): string;
  setMobilId(value: string): TransaksiJualResponse;

  getPenjualId(): string;
  setPenjualId(value: string): TransaksiJualResponse;

  getPembeliId(): string;
  setPembeliId(value: string): TransaksiJualResponse;

  getTotal(): number;
  setTotal(value: number): TransaksiJualResponse;

  getStatus(): string;
  setStatus(value: string): TransaksiJualResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransaksiJualResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TransaksiJualResponse): TransaksiJualResponse.AsObject;
  static serializeBinaryToWriter(message: TransaksiJualResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransaksiJualResponse;
  static deserializeBinaryFromReader(message: TransaksiJualResponse, reader: jspb.BinaryReader): TransaksiJualResponse;
}

export namespace TransaksiJualResponse {
  export type AsObject = {
    id: string,
    mobilId: string,
    penjualId: string,
    pembeliId: string,
    total: number,
    status: string,
  }
}

export class RentMobilRequest extends jspb.Message {
  getMobilId(): string;
  setMobilId(value: string): RentMobilRequest;

  getTanggalMulai(): string;
  setTanggalMulai(value: string): RentMobilRequest;

  getTanggalSelesai(): string;
  setTanggalSelesai(value: string): RentMobilRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RentMobilRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RentMobilRequest): RentMobilRequest.AsObject;
  static serializeBinaryToWriter(message: RentMobilRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RentMobilRequest;
  static deserializeBinaryFromReader(message: RentMobilRequest, reader: jspb.BinaryReader): RentMobilRequest;
}

export namespace RentMobilRequest {
  export type AsObject = {
    mobilId: string,
    tanggalMulai: string,
    tanggalSelesai: string,
  }
}

export class CompleteRentalRequest extends jspb.Message {
  getRentalId(): string;
  setRentalId(value: string): CompleteRentalRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteRentalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteRentalRequest): CompleteRentalRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteRentalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteRentalRequest;
  static deserializeBinaryFromReader(message: CompleteRentalRequest, reader: jspb.BinaryReader): CompleteRentalRequest;
}

export namespace CompleteRentalRequest {
  export type AsObject = {
    rentalId: string,
  }
}

export class TransaksiRentalResponse extends jspb.Message {
  getId(): string;
  setId(value: string): TransaksiRentalResponse;

  getMobilId(): string;
  setMobilId(value: string): TransaksiRentalResponse;

  getPemilikId(): string;
  setPemilikId(value: string): TransaksiRentalResponse;

  getPenyewaId(): string;
  setPenyewaId(value: string): TransaksiRentalResponse;

  getTanggalMulai(): string;
  setTanggalMulai(value: string): TransaksiRentalResponse;

  getTanggalSelesai(): string;
  setTanggalSelesai(value: string): TransaksiRentalResponse;

  getTotal(): number;
  setTotal(value: number): TransaksiRentalResponse;

  getStatus(): string;
  setStatus(value: string): TransaksiRentalResponse;

  getDenda(): number;
  setDenda(value: number): TransaksiRentalResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransaksiRentalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TransaksiRentalResponse): TransaksiRentalResponse.AsObject;
  static serializeBinaryToWriter(message: TransaksiRentalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransaksiRentalResponse;
  static deserializeBinaryFromReader(message: TransaksiRentalResponse, reader: jspb.BinaryReader): TransaksiRentalResponse;
}

export namespace TransaksiRentalResponse {
  export type AsObject = {
    id: string,
    mobilId: string,
    pemilikId: string,
    penyewaId: string,
    tanggalMulai: string,
    tanggalSelesai: string,
    total: number,
    status: string,
    denda: number,
  }
}

export class GetNotificationsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNotificationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetNotificationsRequest): GetNotificationsRequest.AsObject;
  static serializeBinaryToWriter(message: GetNotificationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNotificationsRequest;
  static deserializeBinaryFromReader(message: GetNotificationsRequest, reader: jspb.BinaryReader): GetNotificationsRequest;
}

export namespace GetNotificationsRequest {
  export type AsObject = {
  }
}

export class DashboardSummary extends jspb.Message {
  getTotalMobilAnda(): number;
  setTotalMobilAnda(value: number): DashboardSummary;

  getTransaksiAktif(): number;
  setTransaksiAktif(value: number): DashboardSummary;

  getPendapatanTerakhir(): number;
  setPendapatanTerakhir(value: number): DashboardSummary;

  getNotifikasiBaru(): number;
  setNotifikasiBaru(value: number): DashboardSummary;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DashboardSummary.AsObject;
  static toObject(includeInstance: boolean, msg: DashboardSummary): DashboardSummary.AsObject;
  static serializeBinaryToWriter(message: DashboardSummary, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DashboardSummary;
  static deserializeBinaryFromReader(message: DashboardSummary, reader: jspb.BinaryReader): DashboardSummary;
}

export namespace DashboardSummary {
  export type AsObject = {
    totalMobilAnda: number,
    transaksiAktif: number,
    pendapatanTerakhir: number,
    notifikasiBaru: number,
  }
}

