//
//  AuthViewModel.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-04.
//

import Foundation
import SimpleKeychain

class AuthViewModel: ObservableObject {
    @Published var errorMessage: String?
    @Published var isFetching: Bool
    @Published var isFetchingInitial: Bool
    @Published var user: User?
    @Published var isLoggedIn: Bool
    
    private let userService: UserService
    private let deviceService: DeviceService
    private let deviceID: Int?
    private let keychain: SimpleKeychain
    private let dateFormatter: DateFormatter

    init() {
        self.userService = UserService()
        self.deviceService = DeviceService()
        self.keychain = SimpleKeychain()
        self.dateFormatter = DateFormatter.iso8601Full
        self.isFetching = false
        self.isFetchingInitial = true
        self.user = UserDefaults.standard.getObject(forKey: "user", as: User.self)
        self.isLoggedIn = false
        
        do {
            _ = try self.keychain.string(forKey: "refresh_token")
            let refreshTokenExpiry = dateFormatter.date(from: try self.keychain.string(forKey: "refresh_token_expiry"))
                        
            if let refreshTokenExpiry {
                if refreshTokenExpiry > Date.now {
                    self.isLoggedIn = true
                }
            }
        } catch {
            self.isLoggedIn = false
        }
        
        let storedDeviceID = UserDefaults.standard.integer(forKey: "deviceID")
        self.deviceID = storedDeviceID == 0 ? nil : storedDeviceID
        
        if isLoggedIn && user == nil {
            Task { @MainActor in
                let (res, err) = await userService.fetchUser()
                
                errorMessage = err
                
                guard let res else {
                    isFetchingInitial = false
                    isLoggedIn = false
                    return
                }
                
                user = res
                UserDefaults.standard.setObject(res, forKey: "user")
                isFetchingInitial = false
            }
        } else {
            self.isFetchingInitial = false
        }
    }
}

// MARK: - Login
extension AuthViewModel {
    func login(email: String, password: String) {
        guard validateInput(email: email, password: password) else {
            return
        }
        
        isFetching = true
        
        Task { @MainActor in
            let (res, err) = await userService.login(email: email, password: password, deviceID: deviceID)
        
            errorMessage = err
            
            guard let res else {
                isFetching = false
                return
            }
            
            do {
                try keychain.set(res.accessToken, forKey: "access_token")
                try keychain.set(res.accessTokenExpiresAt, forKey: "access_token_expiry")
                if let refreshToken = res.refreshToken, let refreshTokenExpiresAt = res.refreshTokenExpiresAt {
                    try keychain.set(refreshToken, forKey: "refresh_token")
                    try keychain.set(refreshTokenExpiresAt, forKey: "refresh_token_expiry")
                }
            } catch {
                errorMessage = "An unexpected error has occured. Please try again later."
                isFetching = false
                return
            }
        
            self.user = res.user
            UserDefaults.standard.setObject(res.user, forKey: "user")
            
            if deviceID == nil {
                let (devRes, devErr) = await deviceService.createDevice()
                
                errorMessage = devErr
                
                guard let devRes else {
                    isFetching = false
                    return
                }
                
                do {
                    try keychain.set(devRes.refreshToken, forKey: "refresh_token")
                    try keychain.set(devRes.refreshTokenExpiresAt, forKey: "refresh_token_expiry")
                } catch {
                    errorMessage = "An unexpected error has occured. Please try again later."
                    isFetching = false
                    return
                }
                
                UserDefaults.standard.set(devRes.device.id, forKey: "deviceID")
            }
            
            isFetching = false
            isLoggedIn = true
        }
    }
}

// MARK: - Signup
extension AuthViewModel {
    func signUp(email: String, password: String) {
        guard validateInput(email: email, password: password) else {
            return
        }
        
        isFetching = true
        
        Task { @MainActor in
            let (user, err) = await userService.createUser(email: email, password: password)
            
            errorMessage = err
            guard user != nil else {
                isFetching = false
                return
            }
            
            self.clearUserData()
            self.login(email: email, password: password)
        }
    }
}

// MARK: - Logout
extension AuthViewModel {
    func logout() {
        let refTok = try? self.keychain.string(forKey: "refresh_token")
        
        isFetching = true
        
        Task { @MainActor in
            _ = await userService.logout(refreshToken: refTok ?? "")
            
            self.clearUserData(deleteDevice: true)
            isFetching = false
        }
    }
}

// MARK: - Helpers
extension AuthViewModel {
    private func isValidEmail(_ email: String) -> Bool {
        let emailRegEx = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"

        let emailPred = NSPredicate(format:"SELF MATCHES %@", emailRegEx)
        return emailPred.evaluate(with: email)
    }
    
    private func validateInput(email: String, password: String) -> Bool {
        guard !email.isEmpty, !password.isEmpty else {
            errorMessage = "Please fill in all fields"
            return false
        }
        
        guard isValidEmail(email) else {
            errorMessage = "Please enter a valid email"
            return false
        }
        
        guard password.count >= 6 else {
            errorMessage = "Password must be at least 6 characters"
            return false
        }
        
        return true
    }
    
    private func clearUserData(deleteDevice: Bool = false) {
        UserDefaults.standard.removeObject(forKey: "user")
        if deleteDevice {
            UserDefaults.standard.removeObject(forKey: "deviceID")
        }
        
        try? keychain.deleteItem(forKey: "access_token")
        try? keychain.deleteItem(forKey: "access_token_expiry")
        try? keychain.deleteItem(forKey: "refresh_token")
        try? keychain.deleteItem(forKey: "refresh_token_expiry")
        
        self.user = nil
        self.isLoggedIn = false
    }
}
