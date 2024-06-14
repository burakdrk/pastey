//
//  ProfileView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-11.
//

import SwiftUI
import FluidGradient

struct ProfileView: View {
    @StateObject public var viewModel = ProfileViewModel()
    @EnvironmentObject var auth: AuthViewModel
    
    @State private var showingAlert = false
    @AppStorage("deviceID") var deviceID: Int = -1
    
    var body: some View {
        NavigationStack {
            VStack {
                GradientHeaderView(title: "Profile")
                
                if let user = viewModel.user {
                    VStack {
                        Text(user.email)
                            .font(.title)
                            .bold()
                            .padding(.bottom, 5)
                        
                        if let date = DateFormatter.iso8601Full.date(from: user.createdAt) {
                            Text("Member since: \(date, style: .date)")
                                .font(.subheadline)
                                .foregroundColor(Color.theme.secondaryText)
                        }
                        
                        Text(user.ispremium ? "Premium user" : "Free user")
                            .font(.subheadline)
                            .foregroundColor(user.ispremium ? Color.theme.accent : Color.theme.secondaryText)
                            .padding(.top, 5)
                    }
                    .padding()
                    .frame(maxWidth: .infinity, alignment: .center)
                    
                    DeviceList
                    
                    Spacer()
                    
                    BigButton(text: "Log out", action: {
                        showingAlert = true
                    })
                    .frame(width: 120)
                    .padding()
                    .alert("Are you sure?", isPresented: $showingAlert) {
                        Button("Yes", role: .destructive) {
                            auth.logout()
                        }
                        Button("No", role: .cancel) {}
                    } message: {
                        Text("This will log you out.")
                    }
                } else {
                    Text("No user found")
                        .foregroundColor(Color.theme.secondaryText)
                        .padding()
                        .frame(maxWidth: .infinity, alignment: .center)
                }
                
                Spacer()
            }
            .navigationBarHidden(true)
        }
    }
}

#Preview {
    ProfileView()
        .environmentObject(AuthViewModel())
}

// MARK: - Device List
extension ProfileView {
    private var DeviceList: some View {
        List {
            Section {
                ForEach(viewModel.devices) { device in
                    VStack(alignment: .leading) {
                        
                        Text(device.deviceName)
                            .font(.title3)
                            .bold()
                        
                        HStack {
                            Text(DateFormatter.iso8601Full.date(from: device.createdAt) ?? Date(), style: .date)
                                .font(.subheadline)
                                .foregroundColor(Color.theme.secondaryText)
                            
                            Spacer()
                            
                            if device.id == deviceID {
                                Text("This device")
                                    .font(.subheadline)
                                    .foregroundColor(Color.theme.accent)
                            }
                        }
                    }.deleteDisabled(device.id == deviceID)
                }.onDelete(perform: viewModel.deleteDevice)
            } header: {
                Text("Devices")
                    .font(.title3)
                    .bold()
            }
            .listSectionSeparator(.hidden)
        }
        .listStyle(.plain)
        .refreshable {
            await viewModel.fetchDevices()
        }
    }
}
