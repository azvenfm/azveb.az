
# 🚀 AutoSMM Platform - Cloud Deployment Guide

Sistem artıq bulud mühitinə (Cloud) keçid üçün hazırdır. Lokal problemləri aradan qaldırmaq üçün ən yaxşı yol budur.

## 🛠️ Addım-addım Quraşdırma (Railway.app)

1. **GitHub-a Yüklə:** 
   Kodu GitHub-da yeni bir repository-yə yüklə (Private ola bilər).
   
2. **Railway.app-a Daxil Ol:**
   - [Railway.app](https://railway.app/) saytına daxil ol və GitHub hesabınla giriş et.
   - **"New Project"** $ightarrow$ **"Deploy from GitHub repo"** seç və AutoSMM repository-ni seç.

3. **Bazanı Qur:**
   - Railway-də **"Add Service"** $ightarrow$ **"Database"** $ightarrow$ **"Add PostgreSQL"** seç.
   - Bazanın `DATABASE_URL` dəyərini kopyala.

4. **Backend-i Bağla:**
   - Backend servisinin **"Variables"** bölməsinə get və bunları əlavə et:
     - `DATABASE_URL` = (Kopyaladığın bazanın URL-i)
     - `PORT` = `8080`
     - `JWT_SECRET` = `senin_gizli_acarın_123`

5. **Frontend-i Bağla:**
   - Frontend servisinin **"Variables"** bölməsinə get:
     - `API_URL` = (Backend-in Railway tərəfindən verilən linki)

## 🎉 Nəticə
Sistem avtomatik build olunacaq və sənə `https://autosmm-xxx.up.railway.app` kimi bir link veriləcək.

**Giriş Məlumatları:**
- **İstifadəçi:** `admin`
- **Şifrə:** `Admin123!`
